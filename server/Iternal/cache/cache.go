package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"os"
	"server/Iternal/config"
	"time"
)

type Cache struct {
	db              *redis.Client
	stateExpiration time.Duration
}

func (c Cache) CreateStateCode() (string, error) {
	stateToken := uuid.NewString()
	if err := c.db.Set(context.Background(), stateToken, "true", c.stateExpiration).Err(); err != nil {
		return "", fmt.Errorf("state token set cached err: %v", err)
	}
	return stateToken, nil
}

func (c Cache) GetStateCode(stateToken string) (bool, error) {
	state, err := c.db.Get(context.Background(), stateToken).Result()
	if err != nil {
		return false, fmt.Errorf("state token get cached err: %v", err)
	}
	return state == "true", nil
}

func (c Cache) DeleteStateCode(stateToken string) error {
	if err := c.db.Del(context.Background(), stateToken).Err(); err != nil {
		return fmt.Errorf("state token delete cached err: %v", err)
	}
	return nil
}

func New(cfg config.CacheConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		DB:       cfg.Db,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &Cache{client, cfg.StateExpiration}, nil
}
