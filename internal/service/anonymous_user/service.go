package service

import (
	"context"
	"sync"
	"time"

	"github.com/laiker/chat-server/internal/model"
	"github.com/laiker/chat-server/internal/repository"
	"github.com/laiker/chat-server/internal/service"
)

type anonymousUserService struct {
	repo      repository.AnonymousUserRepository
	cleanupCh chan struct{}
	cleanupWg sync.WaitGroup
}

func NewAnonymousUserService(repo repository.AnonymousUserRepository) service.AnonymousUserService {
	return &anonymousUserService{
		repo:      repo,
		cleanupCh: make(chan struct{}),
	}
}

func (s *anonymousUserService) Create(ctx context.Context, login string) *model.AnonymousUser {
	return s.repo.Create(ctx, login)
}

func (s *anonymousUserService) Get(ctx context.Context, id int64) (*model.AnonymousUser, bool) {
	return s.repo.Get(ctx, id)
}

func (s *anonymousUserService) Remove(ctx context.Context, id int64) bool {
	return s.repo.Remove(ctx, id)
}

func (s *anonymousUserService) StartCleanupRoutine(interval, threshold time.Duration) {
	s.cleanupWg.Add(1)

	go func() {
		defer s.cleanupWg.Done()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.repo.CleanupInactive(threshold)
			case <-s.cleanupCh:
				return
			}
		}
	}()
}

func (s *anonymousUserService) StopCleanupRoutine() {
	close(s.cleanupCh)
	s.cleanupWg.Wait()
}
