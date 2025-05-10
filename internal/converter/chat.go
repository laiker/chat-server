package converter

import (
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

func ToChatFromCreateRequest(chat *chat_v1.CreateRequest) *model.ChatInfo {
	return &model.ChatInfo{
		UsersID: chat.Ids,
		Name:    chat.Name,
		Public:  chat.Public,
	}
}

func ToChatFromConnectRequest(chat *chat_v1.ConnectRequest) *model.ChatConnect {
	return &model.ChatConnect{
		ChatID: chat.ChatId,
		UserID: chat.UserId,
		Login:  chat.Login,
	}
}
