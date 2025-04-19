package chat

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/chat-server/internal/converter"
	"github.com/laiker/chat-server/internal/service"
	"github.com/laiker/chat-server/pkg/chat_v1"
)

type Server struct {
	chat_v1.UnimplementedChatV1Server
	ChatService    service.ChatService
	MessageService service.MessageService
}

func NewServer(chatService service.ChatService, messageService service.MessageService) *Server {
	return &Server{
		ChatService:    chatService,
		MessageService: messageService,
	}
}

func (s *Server) Create(ctx context.Context, request *chat_v1.CreateRequest) (*chat_v1.CreateResponse, error) {

	chatID, err := s.ChatService.Create(ctx, converter.ToChatFromCreateRequest(request))

	if err != nil {
		return nil, err
	}

	return &chat_v1.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *Server) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*empty.Empty, error) {

	_, err := s.MessageService.Create(ctx, converter.ToMessageFromCreateRequest(request))

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, request *chat_v1.DeleteRequest) (*empty.Empty, error) {
	chatId := request.GetId()
	err := s.ChatService.Delete(ctx, chatId)

	if err != nil {
		log.Printf("failed to delete chat: %v", err)
	}

	err = s.MessageService.Delete(ctx, chatId)

	if err != nil {
		log.Printf("failed to delete messages: %v", err)
	}

	return &empty.Empty{}, nil
}
