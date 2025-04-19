package converter

import (
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

func ToMessageFromCreateRequest(message *chat_v1.SendMessageRequest) *model.MessageInfo {
	return &model.MessageInfo{
		ChatID: message.ChatId,
		Value:  message.Text,
		UserID: message.FromUserId,
	}
}
