package converter

import (
	"github.com/laiker/chat-server/pkg/chat_v1"
)

func ToMessageFromCreateRequest(message *chat_v1.SendMessageRequest) *chat_v1.Message {
	return &chat_v1.Message{
		FromUserId: message.Message.FromUserId,
		Text:       message.Message.Text,
		CreatedAt:  message.Message.CreatedAt,
	}
}
