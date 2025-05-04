package model

import "time"

type MessageInfo struct {
	ChatID    int64
	UserID    int64
	Value     string
	CreatedAt time.Time
}
