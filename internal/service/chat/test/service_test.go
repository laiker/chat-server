package test

import (
	"testing"

	"github.com/laiker/chat-server/client/db"
	"github.com/laiker/chat-server/internal/logger/logger"
	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	service "github.com/laiker/chat-server/internal/service/chat"
	. "github.com/ovechkin-dm/mockio/mock"
	"golang.org/x/net/context"
)

type TestDependencies struct {
	ChatRepositoryMock repository.ChatRepository
	txManagerMock      db.TxManager
	loggerMock         logger.DBLoggerInterface
	contextMock        context.Context
}

func SetupServiceTest(t *testing.T) *TestDependencies {
	t.Helper()

	r := Mock[repository.ChatRepository]()
	tx := Mock[db.TxManager]()
	dblogger := Mock[logger.DBLoggerInterface]()

	deps := &TestDependencies{
		ChatRepositoryMock: r,
		txManagerMock:      tx,
		loggerMock:         dblogger,
		contextMock:        context.Background(),
	}

	return deps
}

func Test_serv_Create(t *testing.T) {

	deps := SetupServiceTest(t)

	type fields struct {
		repo      repository.ChatRepository
		txManager db.TxManager
		logger    logger.DBLoggerInterface
	}
	type args struct {
		ctx      context.Context
		ChatInfo *model.ChatInfo
	}

	callback := func(args []any) []any {
		fn := args[1].(func(context.Context) (int64, error))
		id, err := fn(deps.contextMock)
		return []any{id, err}
	}

	When(deps.txManagerMock.ReadCommitted(Any[context.Context](), Any[db.Handler]())).
		ThenReturn(nil).
		ThenAnswer(callback)

	mi := &model.ChatInfo{
		UsersID: []int64{1, 2, 3},
	}

	When(deps.ChatRepositoryMock.Create(deps.contextMock, mi)).ThenReturn(int64(1), nil)

	a := args{
		ctx:      deps.contextMock,
		ChatInfo: mi,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Success Test",
			want:    int64(1),
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.ChatRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewChatService(
				tt.fields.repo, tt.fields.txManager, tt.fields.logger,
			)

			_, err := s.Create(tt.args.ctx, tt.args.ChatInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_serv_Delete(t *testing.T) {
	deps := SetupServiceTest(t)

	type fields struct {
		repo      repository.ChatRepository
		txManager db.TxManager
		logger    logger.DBLoggerInterface
	}
	type args struct {
		ctx context.Context
		id  int64
	}

	a := args{
		ctx: deps.contextMock,
		id:  int64(1),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Success Test",
			wantErr: false,
			args:    a,
			fields: fields{
				repo:      deps.ChatRepositoryMock,
				txManager: deps.txManagerMock,
				logger:    deps.loggerMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewChatService(
				tt.fields.repo, tt.fields.txManager, tt.fields.logger,
			)
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
