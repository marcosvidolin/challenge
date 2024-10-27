package usecase

//go:generate mockgen -source search.go -destination ../../mocks/internal/usecase/search.go -package usecase

import (
	"api/internal/adapter"
	"api/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"
)

type searchCryptor interface {
	Decrypt(user *entity.User) error
}

type searchRepo interface {
	Search(ctx context.Context, queryOpts adapter.QueryOpts) ([]entity.User, error)
}

type searchCache interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, data string, ttl time.Duration) error
}

const (
	searchCacheExp = time.Minute * 15
)

type searchUsecase struct {
	repo    searchRepo
	cache   searchCache
	cryptor searchCryptor
}

func NewSearchUsecase(repo searchRepo, cache searchCache, cryptor searchCryptor) *searchUsecase {
	return &searchUsecase{
		repo:    repo,
		cache:   cache,
		cryptor: cryptor,
	}
}

func (u *searchUsecase) Execute(ctx context.Context, input SearchInput) ([]entity.User, error) {
	key := input.String()
	cached, err := u.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	if cached != nil {
		if err := json.Unmarshal([]byte(*cached), &users); err != nil {
			return nil, err
		}
		decryptedUsers, err := u.decrypt(users)
		if err != nil {
			return nil, err
		}
		return decryptedUsers, nil
	}

	opts := adapter.QueryOpts{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Fields:    input.Fields,
	}

	users, err = u.repo.Search(ctx, opts)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		toCache, err := json.Marshal(users)
		if err != nil {
			slog.Error("search-usecase", slog.Group("Execute", "marshal to cache", err))
			return
		}
		if err := u.cache.Set(ctx, key, string(toCache), searchCacheExp); err != nil {
			slog.Error("search-usecase", slog.Group("Execute", "set user to cache", err))
			return
		}
	}()

	decryptedUsers, err := u.decrypt(users)
	if err != nil {
		return nil, err
	}

	wg.Wait()

	return decryptedUsers, nil
}

func (u *searchUsecase) decrypt(users []entity.User) ([]entity.User, error) {
	dec := make([]entity.User, len(users))
	for i, user := range users {
		if err := u.cryptor.Decrypt(&user); err != nil {
			return nil, err
		}
		dec[i] = user
	}
	return dec, nil
}

type SearchInput struct {
	FirstName string
	LastName  string
	Email     string

	Fields string
}

func (s *SearchInput) SortedFields() []string {
	f := strings.Split(s.Fields, ",")
	sort.Strings(f)
	return f
}

func (s *SearchInput) String() string {
	fields := strings.Join(s.SortedFields(), ",")
	str := fmt.Sprintf("%s-%s-%s-%s", s.FirstName, s.FirstName, s.Email, fields)
	return str
}
