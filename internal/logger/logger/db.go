package logger

import (
	"context"
	_ "context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/chat-server/client/db"
	"github.com/laiker/chat-server/internal/logger"
)

var _ logger.BaseLogger = (*DBLogger)(nil)

type DBLoggerInterface interface {
	logger.BaseLogger
}

type DBLogger struct {
	db db.Client
}

func NewDBLogger(db db.Client) *DBLogger {
	return &DBLogger{db: db}
}

func (l *DBLogger) Log(ctx context.Context, data logger.LogData) error {

	sBuilder := sq.Insert("chat_log").
		Columns("name", "entity_id").
		Values(data.Name, data.EntityID).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return err
	}

	q := db.Query{
		Name:     "log",
		QueryRaw: query,
	}

	_, err = l.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
		return err
	}

	return nil
}
