package chat

import (
	"github.com/laiker/chat-server/client/db"
	log "github.com/laiker/chat-server/internal/logger"
	"github.com/laiker/chat-server/internal/logger/logger"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	"github.com/laiker/chat-server/internal/service"
	"golang.org/x/net/context"
)

type serv struct {
	repo      repository.ChatRepository
	txManager db.TxManager
	logger    logger.DBLoggerInterface
}

func NewChatService(repo repository.ChatRepository, manager db.TxManager, logger logger.DBLoggerInterface) service.ChatService {
	return &serv{repo: repo, txManager: manager, logger: logger}
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

	return id, nil
}

func (s *serv) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
