package model

import (
	"time"
)

type MessageInfo struct {
	ChatID int64
	UserID int64
	Value  string
}

type Message struct {
	MessageInfo
	CreatedAt time.Time
}
