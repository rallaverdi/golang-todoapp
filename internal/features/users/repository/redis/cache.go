package users_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rallaverdi/golang-todoapp/internal/core/domain"
	core_errors "github.com/rallaverdi/golang-todoapp/internal/core/errors"
	"github.com/redis/go-redis/v9"
)

type FilterCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewFilterCache(client *redis.Client, ttl time.Duration) *FilterCache {
	return &FilterCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *FilterCache) FindFilterID(ctx context.Context, filter domain.UsersFilter) (string, error) {
	key, err := paramsKey(filter)
	if err != nil {
		return "", fmt.Errorf("error getting param key: %w", err)
	}

	filterID, err := c.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", core_errors.ErrNotFound
	}

	return filterID, err
}

func (c *FilterCache) SaveResult(
	ctx context.Context,
	filter domain.UsersFilter,
	filterID string,
	users []domain.User,
) error {
	key, err := paramsKey(filter)
	if err != nil {
		return fmt.Errorf("error getting param key: %w", err)
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("error marshalling users: %w", err)
	}

	pipe := c.client.TxPipeline()
	pipe.Set(ctx, key, filterID, c.ttl)
	pipe.Set(ctx, resultKey(filterID), usersJSON, c.ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("error saving result to cache: %w", err)
	}

	return nil
}

func (c *FilterCache) GetUsers(ctx context.Context, filterID string) ([]domain.User, error) {
	usersJSON, err := c.client.Get(ctx, resultKey(filterID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, core_errors.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting users from cache: %w", err)
	}

	var users []domain.User
	if err := json.Unmarshal(usersJSON, &users); err != nil {
		return nil, fmt.Errorf("error unmarshalling users: %w", err)
	}

	return users, nil
}
