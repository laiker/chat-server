package logger

import "context"

type LogData struct {
	Name     string
	EntityID int64
}

type BaseLogger interface {
	Log(ctx context.Context, data LogData) error
}
