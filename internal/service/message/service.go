package chat

import (
	"github.com/laiker/chat-server/client/db"
	log "github.com/laiker/chat-server/internal/logger"
	"github.com/laiker/chat-server/internal/logger/logger"
	"github.com/laiker/chat-server/internal/repository"
	"github.com/laiker/chat-server/internal/service"
	"github.com/laiker/chat-server/pkg/chat_v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serv struct {
	repo           repository.MessageRepository
	txManager      db.TxManager
	logger         logger.DBLoggerInterface
	channelManager *service.ChatChannelManager
}

func NewMessageService(repo repository.MessageRepository, manager db.TxManager, logger logger.DBLoggerInterface, channelManager *service.ChatChannelManager) service.MessageService {
	return &serv{repo: repo, txManager: manager, logger: logger, channelManager: channelManager}
}

func (s *serv) Create(ctx context.Context, chatId int64, messageInfo *chat_v1.Message) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.repo.Create(ctx, chatId, messageInfo)

		if errTx != nil {
			return errTx
		}

		logData := log.LogData{
			Name:     "create message",
			EntityID: id,
		}

		errTx = s.logger.Log(ctx, logData)

		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	s.channelManager.M.RLock()
	chatChan, ok := s.channelManager.Channels[chatId]
	defer s.channelManager.M.RUnlock()

	if !ok {
		return 0, status.Errorf(codes.NotFound, "chat not found")
	}

	chatChan <- messageInfo

	return id, nil
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
