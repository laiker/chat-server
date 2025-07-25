package service

import (
	"context"
	"sync"
	"time"

	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

type ChatService interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Delete(ctx context.Context, id int64) error
	Connect(info model.ChatConnect, stream chat_v1.ChatV1_ConnectServer) error
	InitializeConnection(info model.ChatConnect, stream chat_v1.ChatV1_ConnectServer) error
	GetUserChats(ctx context.Context, id int64) ([]model.Chat, error)
}

type MessageService interface {
	Create(ctx context.Context, chatId int64, info *chat_v1.Message) (int64, error)
	Delete(ctx context.Context, id int64) error
	GetHistory(ctx context.Context, id int64) ([]*chat_v1.Message, error)
	SendToChannel(ctx context.Context, chatId int64, info *chat_v1.Message) (bool, error)
}

type AnonymousUserService interface {
	Create(ctx context.Context, login string) *model.AnonymousUser
	Get(ctx context.Context, id int64) (*model.AnonymousUser, bool)
	Remove(ctx context.Context, id int64) bool
	StartCleanupRoutine(interval, threshold time.Duration)
	StopCleanupRoutine()
}

type ChatStreamManager struct {
	Streams map[int64]*model.ChatStream
	M       sync.RWMutex
}

type ChatChannelManager struct {
	Channels map[int64]chan *chat_v1.Message
	M        sync.RWMutex
}
