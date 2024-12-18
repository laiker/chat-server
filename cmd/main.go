package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	chat_v1.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {
	fmt.Printf("%+v", request)
	return &chat_v1.CreateResponse{}, nil
}

func (s *server) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*empty.Empty, error) {
	fmt.Printf("%+v", request)
	return &empty.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*empty.Empty, error) {
	fmt.Printf("%+v", request)
	return &empty.Empty{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50053")

	if err != nil {
		panic(err)
	}

	g := grpc.NewServer()
	reflection.Register(g)
	chat_v1.RegisterChatV1Server(g, &server{})

	log.Printf("server listening at %v", listener.Addr())

	if err = g.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
