package model

import (
	"sync"
	"time"

	"github.com/laiker/chat-server/pkg/chat_v1"
)

type ChatInfo struct {
	UsersID []int64
	Name    string
	Public  bool
}

type Chat struct {
	Id        int64     `db:"id"`
	UsersID   []int64   `db:"users_id"`
	CreatedAt time.Time `db:"created_at"`
}

type ChatConnect struct {
	ChatID int64
	UserID int64
}

type ChatStream struct {
	Streams map[int64]chat_v1.ChatV1_ConnectServer
	M       sync.RWMutex
}
