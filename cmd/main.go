package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/laiker/chat-server/internal/config"
	"github.com/laiker/chat-server/internal/config/env"
	"github.com/laiker/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	chat_v1.UnimplementedChatV1Server
	db *pgxpool.Pool
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func (s *server) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	sBuilder := sq.Insert("chat").
		Columns("created_at").
		Values(time.Now()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var chatID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&chatID)

	if err != nil {
		log.Fatalf("failed to create chat: %v", err)
	}

	sBuilder = sq.Insert("chat_user").
		Columns("chat_id", "user_id").
		Values(chatID, 1).
		PlaceholderFormat(sq.Dollar)

	query, args, err = sBuilder.ToSql()

	s.db.QueryRow(ctx, query, args...)

	return &chat_v1.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*empty.Empty, error) {
	sBuilder := sq.Insert("message").
		Columns("chat_id", "user_id", "message").
		Values(request.ChatId, request.From, request.Text).
		PlaceholderFormat(sq.Dollar)

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, errm := s.db.Exec(ctx, query, args...)

	if errm != nil {
		log.Fatalf("failed to create message: %v", errm)
	}

	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*empty.Empty, error) {
	sBuilder := sq.Delete("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": request.Id})

	query, args, err := sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	sBuilder = sq.Delete("chat_user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"chat_id": request.Id})

	query, args, err = sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.db.Exec(ctx, query, args...)

	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	sBuilder = sq.Delete("message").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"chat_id": request.Id})

	query, args, err = sBuilder.ToSql()

	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	_, err = s.db.Exec(ctx, query, args...)

	if err != nil {
		log.Fatalf("failed to delete chat: %v", err)
	}

	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()

	errConfig := config.Load(configPath)
	ctx := context.Background()

	if errConfig != nil {
		log.Fatalf("failed to load config: %v", errConfig)
	}

	gConfig, err := env.NewGRPCConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	listener, err := net.Listen("tcp", gConfig.Address())

	if err != nil {
		log.Fatalf("file to start server: %v", err)
	}

	pgConfig, errc := env.NewPGConfig()

	if errc != nil {
		log.Fatalf("failed to load config: %v", errc)
	}

	p, errp := pgxpool.Connect(ctx, pgConfig.DSN())

	if errp != nil {
		log.Fatalf("failed to connect: %v", errp)
	}

	g := grpc.NewServer()
	reflection.Register(g)

	chat_v1.RegisterChatV1Server(g, &server{db: p})

	log.Printf("server listening at %v", listener.Addr())

	if err = g.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if errp != nil {
		log.Fatalf("failed to connect: %v", errp)
	}

	defer p.Close()
}
