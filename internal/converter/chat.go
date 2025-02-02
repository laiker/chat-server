package converter

import (
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

func ToChatFromCreateRequest(chat *chat_v1.CreateRequest) *model.ChatInfo {
	values := []int64{1, 2, 3, 4}
	return &model.ChatInfo{
		UsersID: values,
	}
}
