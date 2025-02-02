package service

import (
	"context"

	"github.com/laiker/chat-server/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, info *model.ChatInfo) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MessageService interface {
	Create(ctx context.Context, info *model.MessageInfo) (int64, error)
	Delete(ctx context.Context, id int64) error
}
