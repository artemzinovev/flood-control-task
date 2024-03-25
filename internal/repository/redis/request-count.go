package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"task/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type RequestCountRepository struct {
	client *redis.Client

	ttl time.Duration
}

func NewRequestCountRepository(addr string, password string, db int, TTL time.Duration) (*RequestCountRepository, error) {
	const op = "redis.NewRequestCountRepository"

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RequestCountRepository{
		client: client,
		ttl:    TTL,
	}, nil
}

func (r *RequestCountRepository) Set(ctx context.Context, userID int64, requestCount int) error {
	const op = "redis.RequestCountRepository.Set"

	err := r.client.Set(ctx, strconv.FormatInt(userID, 10), requestCount, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RequestCountRepository) Get(ctx context.Context, userID int64) (int, error) {
	const op = "redis.RequestCountRepository.Get"

	result, err := r.client.Get(ctx, strconv.FormatInt(userID, 10)).Int()
	if errors.Is(err, redis.Nil) {
		return -1, repository.ErrKeyNotExist
	}

	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
