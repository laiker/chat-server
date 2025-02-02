package model

import (
	"time"
)

type ChatInfo struct {
	UsersID []int64
}

type Chat struct {
	Id        int64
	UsersID   []int64
	CreatedAt time.Time
}
