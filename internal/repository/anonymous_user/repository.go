package repository

import (
	"context"
	"sync"
	"time"

	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
)

type anonymousUserRepo struct {
	users map[int64]*model.AnonymousUser
	mu    sync.RWMutex
}

func NewAnonymousUserRepository(ctx context.Context) repository.AnonymousUserRepository {
	return &anonymousUserRepo{
		users: make(map[int64]*model.AnonymousUser),
	}
}

func (r *anonymousUserRepo) Create(ctx context.Context, login string) *model.AnonymousUser {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := model.NewAnonymousUser(login)
	r.users[user.GetID()] = user

	return user
}

func (r *anonymousUserRepo) Get(ctx context.Context, id int64) (*model.AnonymousUser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	return user, exists
}

func (r *anonymousUserRepo) GetAll(ctx context.Context) []*model.AnonymousUser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.AnonymousUser, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}

func (r *anonymousUserRepo) Remove(ctx context.Context, id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[id]
	if exists {
		delete(r.users, id)
		return true
	}

	return false
}

func (r *anonymousUserRepo) CleanupInactive(threshold time.Duration) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-threshold)
	removed := 0

	for id, user := range r.users {
		if user.GetCreatedAt().Before(cutoff) {
			delete(r.users, id)
			removed++
		}
	}

	return removed
}
