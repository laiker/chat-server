package repository

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/chat-server/client/db"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	"golang.org/x/net/context"
)

const (
	tableName = "message"

	chatIdColumn    = "chat_id"
	userIdColumn    = "user_id"
	messageColumn   = "message"
	createdAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

func NewMessageRepository(db db.Client) repository.MessageRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, messageInfo *model.MessageInfo) (int64, error) {

	sBuilder := sq.Insert(tableName).
		Columns(chatIdColumn, userIdColumn, messageColumn, createdAtColumn).
		Values(messageInfo.ChatID, messageInfo.UserID, messageInfo.Value, time.Now()).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	var userID int64

	q := db.Query{
		Name:     "message.create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)

	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
	}

	return userID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {

	sBuilder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIdColumn: id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "message.create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to delete chat: %v", err)
	}

	return nil
}
