package chat

import (
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
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
	nameColumn        = "name"
	publicColumn      = "public"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chatInfo *model.ChatInfo) (int64, error) {

	sBuilder := sq.Insert(chatTableName).
		Columns(createdAtColumn, nameColumn, publicColumn).
		Values(time.Now(), chatInfo.Name, chatInfo.Public).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	fmt.Println(query, args)

	if err != nil {
		log.Printf("failed to build query: %v\n", err)
	}

	var chatID int64

	q := db.Query{
		Name:     "chat.create",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)

	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
	}

	for i := 0; i < len(chatInfo.UsersID); i++ {

		sBuilder = sq.Insert(chatUserTableName).
			Columns(chatIDColumn, userIDColumn).
			Values(chatID, chatInfo.UsersID[i]).
			PlaceholderFormat(sq.Dollar)

		query, args, err = sBuilder.ToSql()

		fmt.Println(query, args)

		if err != nil {
			log.Printf("failed to build query: %v\n", err)
		}

		q = db.Query{
			Name:     "chat_user.create",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)

		if err != nil {
			log.Printf("failed to insert user: %v\n", err)
		}
	}
	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {

	sBuilder := sq.Delete(chatTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat.delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to delete chat: %v", err)
	}

	sBuilder = sq.Delete(chatUserTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{chatIDColumn: id})

	query, args, err = sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	q = db.Query{
		Name:     "chat_user.delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to delete chat: %v", err)
	}

	return nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Chat, error) {

	sBuilder := sq.Select(idColumn, "array_agg(cu.user_id) AS users_id", createdAtColumn).
		From(chatTableName).
		LeftJoin("chat_user cu ON cu.chat_id = chat.id").
		Where(sq.Eq{idColumn: id}).
		GroupBy(idColumn).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat.get",
		QueryRaw: query,
	}

	var chat model.Chat

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chat.Id, &chat.UsersID, &chat.CreatedAt)

	if err != nil {
		fmt.Printf("failed to find chat: %v", err)
	}

	if chat.Id <= 0 {
		return nil, fmt.Errorf("chat not found")
	}

	return &chat, nil
}

func (r *repo) GetUserChats(ctx context.Context, userID int64) ([]model.Chat, error) {

	sBuilder := sq.Select(idColumn, "array_agg(cu.user_id) AS users_id", createdAtColumn, nameColumn, publicColumn).
		From(chatTableName).
		LeftJoin("chat_user cu ON cu.chat_id = chat.id").
		Where(sq.Or{
			sq.Eq{userIDColumn: userID},
			sq.Eq{publicColumn: true},
		}).
		GroupBy(idColumn).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, err
	}

	q := db.Query{
		Name:     "chat.get_user_chats",
		QueryRaw: query,
	}

	var chats []model.Chat

	rows, err := r.db.DB().QueryContext(ctx, q, args...)

	if err != nil {
		log.Printf("failed to find chats: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var chat model.Chat
		var usersIDSlice pgtype.Int4Array
		err = rows.Scan(&chat.Id, &usersIDSlice, &chat.CreatedAt, &chat.Name, &chat.Public)

		if err != nil {
			log.Printf("failed to scan chat: %v", err)
			return nil, err
		}

		for _, elem := range usersIDSlice.Elements {
			chat.UsersID = append(chat.UsersID, int64(elem.Int))
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
