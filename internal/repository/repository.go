package repository

import (
	"context"

	"github.com/laiker/chat-server/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	Create(ctx context.Context, info *model.MessageInfo) (int64, error)
	Delete(ctx context.Context, id int64) error
}
