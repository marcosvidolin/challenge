package usecase

//go:generate mockgen -source get_by_id.go -destination ../../mocks/internal/usecase/mock_get_by_id.go -package usecase

import (
	"api/internal/entity"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type getByIDCryptor interface {
	Decrypt(user *entity.User) error
}

type getByIDRepo interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
}

type getByIdCache interface {
	Get(ctx context.Context, id string) (*string, error)
	Set(ctx context.Context, key string, data string, ttl time.Duration) error
}

const (
	getByIDCacheExp = time.Minute * 15
)

type getByIDUsecase struct {
	repo    getByIDRepo
	cache   getByIdCache
	cryptor getByIDCryptor
}

func NewGetByIDUsecase(repo getByIDRepo, cache getByIdCache, cryptor getByIDCryptor) *getByIDUsecase {
	return &getByIDUsecase{
		repo:    repo,
		cache:   cache,
		cryptor: cryptor,
	}
}

func (u *getByIDUsecase) Execute(ctx context.Context, id string) (*entity.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, entity.NewBusinessError("invalid empty id")
	}

	cached, err := u.cache.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	var user entity.User

	if cached != nil {
		if err := json.Unmarshal([]byte(*cached), &user); err != nil {
			return nil, err
		}
		if err := u.cryptor.Decrypt(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	userData, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(user entity.User) {
		defer wg.Done()
		toCache, err := json.Marshal(user)
		if err != nil {
			slog.Error("getByID-usecase", slog.Group("Execute", "marshal to cache", err))
			return
		}
		if err = u.cache.Set(ctx, id, string(toCache), getByIDCacheExp); err != nil {
			slog.Error("getByID-usecase", slog.Group("Execute", "set user to cache", err))
			return
		}
	}(*userData)

	if err := u.cryptor.Decrypt(userData); err != nil {
		return nil, err
	}

	wg.Wait()

	return userData, nil
}
