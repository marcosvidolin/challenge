package usecase

import (
	"api/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type upsertQueue interface {
	Consume(ch chan<- string) error
}

type upsertRepo interface {
	Upsert(ctx context.Context, user entity.User) error
}

type upsertCache interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, data string, ttl time.Duration) error
}

const (
	upsertCacheExp = time.Minute * 15
)

type upsertUsecase struct {
	queue    upsertQueue
	userRepo upsertRepo
	cache    upsertCache
}

func NewUpsertUsecase(queue upsertQueue, repo upsertRepo, cache upsertCache) *upsertUsecase {
	return &upsertUsecase{
		queue:    queue,
		userRepo: repo,
		cache:    cache,
	}
}

func (u *upsertUsecase) Execute(ctx context.Context) error {
	messageChannel := make(chan string)
	if err := u.queue.Consume(messageChannel); err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}
	defer close(messageChannel)

	var wg sync.WaitGroup
	for msg := range messageChannel {
		var users []entity.User
		if err := json.Unmarshal([]byte(msg), &users); err != nil {
			slog.Error("upsert-usecase", slog.Group("Execute", "unmarshal", err))
			continue
		}

		for _, user := range users {
			wg.Add(1)
			if err := u.userRepo.Upsert(ctx, user); err != nil {
				slog.Error("upsert-usecase", slog.Group("Execute", "upsert", err))
				continue
			}
			go func(user entity.User) {
				defer wg.Done()
				toCache, err := json.Marshal(user)
				if err != nil {
					slog.Error("upsert-usecase", slog.Group("Execute", "cache marshal", err))
					return
				}
				key := user.Email
				if err := u.cache.Set(ctx, key, string(toCache), upsertCacheExp); err != nil {
					slog.Error("upsert-usecase", slog.Group("Execute", "cache set", err))
					return
				}
			}(user)
		}
	}

	wg.Wait()

	return nil
}
