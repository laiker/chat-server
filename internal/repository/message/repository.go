package repository

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/laiker/chat-server/client/db"
	"github.com/laiker/chat-server/internal/repository"
	"github.com/laiker/chat-server/pkg/chat_v1"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	tableName = "message"

	chatIdColumn    = "chat_id"
	userIdColumn    = "user_id"
	messageColumn   = "message"
	createdAtColumn = "created_at"
	loginColumn     = "login"
)

type repo struct {
	db db.Client
}

func NewMessageRepository(db db.Client) repository.MessageRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatId int64, messageInfo *chat_v1.Message) (int64, error) {
	log.Println(messageInfo.UserLogin)
	sBuilder := sq.Insert(tableName).
		Columns(chatIdColumn, userIdColumn, messageColumn, createdAtColumn, loginColumn).
		Values(chatId, messageInfo.GetUserId(), messageInfo.GetText(), messageInfo.GetCreatedAt().AsTime(), messageInfo.UserLogin).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	var messageID int64

	q := db.Query{
		Name:     "message.create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&messageID)

	if err != nil {
		log.Printf("failed to insert message: %v\n", err)
	}

	return messageID, nil
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

func (r *repo) GetHistory(ctx context.Context, chatId int64, limit int64) ([]*chat_v1.Message, error) {

	sBuilder := sq.Select(messageColumn, createdAtColumn, userIdColumn, loginColumn).
		From(tableName).
		Where(sq.Eq{chatIdColumn: chatId}).
		OrderBy(createdAtColumn + " ASC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "message.get",
		QueryRaw: query,
	}

	var messages []*chat_v1.Message

	rows, err := r.db.DB().QueryContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to find messages: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var message chat_v1.Message
		var messageTime time.Time

		err = rows.Scan(&message.Text, &messageTime, &message.UserId, &message.UserLogin)

		if err != nil {
			log.Printf("failed to scan message: %v", err)
		}

		message.CreatedAt = timestamppb.New(messageTime)

		messages = append(messages, &message)
	}

	return messages, nil
}
