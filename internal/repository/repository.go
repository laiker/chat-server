package repository

import (
	"context"

	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

type ChatRepository interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Chat, error)
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	Create(ctx context.Context, chatId int64, info *chat_v1.Message) (int64, error)
	Delete(ctx context.Context, id int64) error
}
