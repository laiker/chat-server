package test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/laiker/chat-server/internal/api/chat"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/service"
	"github.com/laiker/chat-server/pkg/chat_v1"
	. "github.com/ovechkin-dm/mockio/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fields struct {
	UnimplementedChatV1Server chat_v1.UnimplementedChatV1Server
	ChatService               service.ChatService
	MessageService            service.MessageService
}

type TestDependencies struct {
	ChatServiceMock service.ChatService
	context         context.Context
}

func SetupApiTest(t *testing.T) *TestDependencies {
	t.Helper()

	m := Mock[service.ChatService]()

	deps := &TestDependencies{
		ChatServiceMock: m,
		context:         context.Background(),
	}

	return deps
}

func TestServer_Delete(t *testing.T) {

	type args struct {
		ctx     context.Context
		request *chat_v1.DeleteRequest
	}

	c := Mock[service.ChatService]()
	m := Mock[service.MessageService]()

	a := args{
		ctx: context.Background(),
		request: &chat_v1.DeleteRequest{
			Id: 1,
		},
	}

	When(c.Delete(a.ctx, 1)).ThenReturn(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    &empty.Empty{},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedChatV1Server: chat_v1.UnimplementedChatV1Server{},
				ChatService:               c,
				MessageService:            m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chat.Server{
				UnimplementedChatV1Server: tt.fields.UnimplementedChatV1Server,
				ChatService:               tt.fields.ChatService,
				MessageService:            tt.fields.MessageService,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_SendMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *chat_v1.SendMessageRequest
	}

	c := Mock[service.ChatService]()
	m := Mock[service.MessageService]()

	a := args{
		ctx: context.Background(),
		request: &chat_v1.SendMessageRequest{
			From:      "1",
			Text:      "message",
			ChatId:    1,
			Timestamp: timestamppb.New(time.Now()),
		},
	}

	mi := &model.MessageInfo{
		ChatID: 1,
		UserID: 1,
		Value:  "message",
	}

	WhenDouble(m.Create(a.ctx, mi)).ThenReturn(1, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    &empty.Empty{},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedChatV1Server: chat_v1.UnimplementedChatV1Server{},
				ChatService:               c,
				MessageService:            m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chat.Server{
				UnimplementedChatV1Server: tt.fields.UnimplementedChatV1Server,
				ChatService:               tt.fields.ChatService,
				MessageService:            tt.fields.MessageService,
			}
			got, err := s.SendMessage(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Create(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *chat_v1.CreateRequest
	}

	c := Mock[service.ChatService]()
	m := Mock[service.MessageService]()

	a := args{
		ctx: context.Background(),
		request: &chat_v1.CreateRequest{
			Usernames: []string{"1", "2", "3", "4"},
		},
	}

	mci := &model.ChatInfo{
		UsersID: []int64{1, 2, 3, 4},
	}

	WhenDouble(c.Create(a.ctx, mci)).ThenReturn(int64(1), nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *chat_v1.CreateResponse
		wantErr bool
	}{
		{
			name: "Success Test",
			want: &chat_v1.CreateResponse{
				Id: int64(1),
			},
			wantErr: false,
			args:    a,
			fields: fields{
				UnimplementedChatV1Server: chat_v1.UnimplementedChatV1Server{},
				ChatService:               c,
				MessageService:            m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chat.Server{
				UnimplementedChatV1Server: tt.fields.UnimplementedChatV1Server,
				ChatService:               tt.fields.ChatService,
				MessageService:            tt.fields.MessageService,
			}

			got, err := s.Create(tt.args.ctx, tt.args.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
