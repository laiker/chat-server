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
	ChatService      service.ChatService
	MessageService   service.MessageService
	AnonymousService service.AnonymousUserService
}

func NewServer(chatService service.ChatService, messageService service.MessageService, anonymouseService service.AnonymousUserService) *Server {
	return &Server{
		ChatService:      chatService,
		MessageService:   messageService,
		AnonymousService: anonymouseService,
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

func (s *Server) Connect(request *chat_v1.ConnectRequest, stream chat_v1.ChatV1_ConnectServer) error {

	err := s.ChatService.Connect(*converter.ToChatFromConnectRequest(request), stream)

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) SendMessage(ctx context.Context, request *chat_v1.SendMessageRequest) (*empty.Empty, error) {

	_, err := s.MessageService.Create(ctx, request.ChatId, request.Message)

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

func (s *Server) CreateAnonymousUser(ctx context.Context, request *chat_v1.CreateAnonymousUserRequest) (*chat_v1.CreateAnonymousUserResponse, error) {
	login := request.GetLogin()

	anonUser := s.AnonymousService.Create(ctx, login)

	return &chat_v1.CreateAnonymousUserResponse{
		UserId: anonUser.GetID(),
		Login:  anonUser.GetLogin(),
	}, nil
}
