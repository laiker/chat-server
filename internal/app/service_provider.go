package app

import (
	"context"
	"log"
	"sync"

	"github.com/laiker/chat-server/client/db"
	"github.com/laiker/chat-server/client/db/pg"
	"github.com/laiker/chat-server/client/db/transaction"
	api "github.com/laiker/chat-server/internal/api/chat"
	"github.com/laiker/chat-server/internal/config"
	"github.com/laiker/chat-server/internal/config/env"
	"github.com/laiker/chat-server/internal/logger/logger"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	repoAnonUser "github.com/laiker/chat-server/internal/repository/anonymous_user"
	repoChat "github.com/laiker/chat-server/internal/repository/chat"
	repoMessage "github.com/laiker/chat-server/internal/repository/message"
	"github.com/laiker/chat-server/internal/service"
	serviceAnonUser "github.com/laiker/chat-server/internal/service/anonymous_user"
	servChat "github.com/laiker/chat-server/internal/service/chat"
	servMessage "github.com/laiker/chat-server/internal/service/message"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

type ServiceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	chatRepository      repository.ChatRepository
	messageRepository   repository.MessageRepository
	chatService         service.ChatService
	messageService      service.MessageService
	chatApi             *api.Server
	db                  db.Client
	txManager           db.TxManager
	dbLogger            *logger.DBLogger
	chatStream          *service.ChatStreamManager
	chatChannel         *service.ChatChannelManager
	anonymousService    service.AnonymousUserService
	anonymousRepository repository.AnonymousUserRepository
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) PGConfig() config.PGConfig {

	if s.pgConfig == nil {
		pgConfig, err := env.NewPGConfig()

		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.pgConfig = pgConfig
	}

	return s.pgConfig
}

func (s *ServiceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {

		gConfig, err := env.NewGRPCConfig()

		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		s.grpcConfig = gConfig

	}

	return s.grpcConfig
}

func (s *ServiceProvider) DB(ctx context.Context) db.Client {
	if s.db == nil {
		p, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect: %v", err)
		}

		s.db = p
	}
	return s.db
}

func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DB(ctx).DB())
	}

	return s.txManager
}

func (s *ServiceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {

	if s.chatRepository == nil {
		r := repoChat.NewChatRepository(s.DB(ctx))
		s.chatRepository = r
	}

	return s.chatRepository
}

func (s *ServiceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {

	if s.messageRepository == nil {
		r := repoMessage.NewMessageRepository(s.DB(ctx))
		s.messageRepository = r
	}

	return s.messageRepository
}

func (s *ServiceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		r := servChat.NewChatService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
			s.DBLogger(ctx),
			s.ChatStreamManager(),
			s.ChatChannelManager(),
		)
		s.chatService = r
	}

	return s.chatService
}

func (s *ServiceProvider) MessageService(ctx context.Context) service.MessageService {

	if s.messageService == nil {
		r := servMessage.NewMessageService(s.MessageRepository(ctx), s.TxManager(ctx), s.DBLogger(ctx), s.ChatChannelManager())
		s.messageService = r
	}

	return s.messageService
}

func (s *ServiceProvider) ChatApi(ctx context.Context) *api.Server {
	if s.chatApi == nil {
		a := api.NewServer(s.ChatService(ctx), s.MessageService(ctx), s.AnonymouseUserService(ctx))
		s.chatApi = a
	}

	return s.chatApi
}

func (s *ServiceProvider) DBLogger(ctx context.Context) logger.DBLoggerInterface {
	if s.dbLogger == nil {
		l := logger.NewDBLogger(s.DB(ctx))
		s.dbLogger = l
	}

	return s.dbLogger
}

func (s *ServiceProvider) ChatStreamManager() *service.ChatStreamManager {
	if s.chatStream == nil {
		s.chatStream = &service.ChatStreamManager{
			Streams: make(map[int64]*model.ChatStream),
			M:       sync.RWMutex{},
		}
	}

	return s.chatStream
}

func (s *ServiceProvider) ChatChannelManager() *service.ChatChannelManager {
	if s.chatChannel == nil {
		s.chatChannel = &service.ChatChannelManager{
			Channels: make(map[int64]chan *chat_v1.Message),
			M:        sync.RWMutex{},
		}
	}

	return s.chatChannel
}

func (s *ServiceProvider) AnonymouseUserRepository(ctx context.Context) repository.AnonymousUserRepository {
	if s.anonymousRepository == nil {
		r := repoAnonUser.NewAnonymousUserRepository(ctx)
		s.anonymousRepository = r
	}

	return s.anonymousRepository
}

func (s *ServiceProvider) AnonymouseUserService(ctx context.Context) service.AnonymousUserService {
	if s.anonymousService == nil {
		r := serviceAnonUser.NewAnonymousUserService(s.AnonymouseUserRepository(ctx))
		s.anonymousService = r
	}

	return s.anonymousService
}
