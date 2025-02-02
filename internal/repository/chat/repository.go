package chat

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
	chatTableName     = "chat"
	chatUserTableName = "chat_user"
	idColumn          = "id"
	chatIDColumn      = "chat_id"
	userIDColumn      = "user_id"
	createdAtColumn   = "created_at"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatInfo *model.ChatInfo) (int64, error) {

	sBuilder := sq.Insert(chatTableName).
		Columns(createdAtColumn).
		Values(time.Now()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Println("failed to build query: %v", err)
	}

	var chatID int64

	q := db.Query{
		Name:     "chat.create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)

	if err != nil {
		log.Println("failed to insert user: %v", err)
	}

	for i := 0; i < len(chatInfo.UsersID); i++ {
		sBuilder = sq.Insert(chatUserTableName).
			Columns(chatIDColumn, userIDColumn).
			Values(chatID, chatInfo.UsersID).
			PlaceholderFormat(sq.Dollar)

		query, args, err = sBuilder.ToSql()

		fmt.Println(query, args)

		if err != nil {
			log.Println("failed to build query: %v", err)
		}

		q = db.Query{
			Name:     "chat_user.create",
			QueryRaw: query,
		}

		err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	}
	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {

	sBuilder := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat.delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	sBuilder = sq.Delete(chatUserTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err = sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q = db.Query{
		Name:     "chat_user.delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	return nil
}
