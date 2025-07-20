package chat

import (
	"fmt"

	"github.com/laiker/chat-server/client/db"
	log "github.com/laiker/chat-server/internal/logger"
	"github.com/laiker/chat-server/internal/logger/logger"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	"github.com/laiker/chat-server/internal/service"
	"github.com/laiker/chat-server/pkg/chat_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serv struct {
	repo      repository.ChatRepository
	txManager db.TxManager
	logger    logger.DBLoggerInterface

	streamManager  *service.ChatStreamManager
	channelManager *service.ChatChannelManager
}

func NewChatService(
	repo repository.ChatRepository,
	manager db.TxManager,
	logger logger.DBLoggerInterface,
	streamManager *service.ChatStreamManager,
	channelManager *service.ChatChannelManager) service.ChatService {

	return &serv{
		repo:           repo,
		txManager:      manager,
		logger:         logger,
		streamManager:  streamManager,
		channelManager: channelManager,
	}
}

func (s *serv) Create(ctx context.Context, chatInfo *model.ChatInfo) (int64, error) {

	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.repo.Create(ctx, chatInfo)
		if errTx != nil {
			return errTx
		}

		logData := log.LogData{
			Name:     "create chat",
			EntityID: id,
		}

		errTx = s.logger.Log(ctx, logData)

		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return id, err
	}

	s.channelManager.M.Lock()
	s.channelManager.Channels[id] = make(chan *chat_v1.Message, 100)
	s.channelManager.M.Unlock()

	return id, nil
}

func (s *serv) GetUserChats(ctx context.Context, userId int64) ([]model.Chat, error) {
	return s.repo.GetUserChats(ctx, userId)
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *serv) InitializeConnection(connect model.ChatConnect, stream chat_v1.ChatV1_ConnectServer) error {
	fmt.Printf("connect: %+v\n", connect)

	s.channelManager.M.RLock()
	_, ok := s.channelManager.Channels[connect.ChatID]
	s.channelManager.M.RUnlock()

	if !ok {
		_, err := s.repo.Get(stream.Context(), connect.ChatID)
		if err != nil {
			return status.Errorf(codes.NotFound, err.Error())
		}

		s.channelManager.M.Lock()
		s.channelManager.Channels[connect.ChatID] = make(chan *chat_v1.Message, 100)
		s.channelManager.M.Unlock()
	}

	// Регистрируем stream
	s.streamManager.M.Lock()
	if _, okChat := s.streamManager.Streams[connect.ChatID]; !okChat {
		s.streamManager.Streams[connect.ChatID] = &model.ChatStream{
			Streams: make(map[int64]chat_v1.ChatV1_ConnectServer),
		}
	}
	chatStream := s.streamManager.Streams[connect.ChatID]
	s.streamManager.M.Unlock()

	chatStream.M.Lock()
	chatStream.Streams[connect.UserID] = stream
	chatStream.M.Unlock()

	return nil
}

func (s *serv) Connect(connect model.ChatConnect, stream chat_v1.ChatV1_ConnectServer) error {
	s.channelManager.M.RLock()
	chatChan := s.channelManager.Channels[connect.ChatID]
	s.channelManager.M.RUnlock()

	s.streamManager.M.RLock()
	chatStream := s.streamManager.Streams[connect.ChatID]
	s.streamManager.M.RUnlock()

	// Отправляем сообщение о подключении
	chatChan <- &chat_v1.Message{
		UserId: 0,
		Text:   fmt.Sprintf("Пользователь %s подключился к чату", connect.Login),
	}

	// Бесконечный цикл для новых сообщений
	for {
		select {
		case msg, okCh := <-chatChan:
			if !okCh {
				return nil
			}

			chatStream.M.RLock()
			for _, st := range chatStream.Streams {
				if err := st.Send(msg); err != nil {
					chatStream.M.RUnlock()
					return err
				}
			}
			chatStream.M.RUnlock()

		case <-stream.Context().Done():
			chatStream.M.Lock()
			delete(chatStream.Streams, connect.UserID)
			chatStream.M.Unlock()

			chatChan <- &chat_v1.Message{
				UserId: 0,
				Text:   fmt.Sprintf("Пользователь %s вышел из чата", connect.Login),
			}
			return nil
		}
	}
}
