package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"server/internal/config"
	"time"
)

type Cache struct {
	Db                           *redis.Client
	StateExpiration              time.Duration
	EmailConfirmedCodeExpiration time.Duration
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

	return &Cache{client, cfg.StateExpiration, cfg.EmailConfirmedCodeExpiration}, nil
}

//func (c Cache) CreateEmailConfirmedCode(email string, code string) (string, error) {
//	if err := c.db.Set(context.Background(), email, code, c.emailConfirmedCodeExpiration).Err(); err != nil {
//		return "", fmt.Errorf("email code set cached err: %v", err)
//	}
//	return code, nil
//}
//
//func (c Cache) GetEmailConfirmedCode(email string) (string, error) {
//	code, err := c.db.Get(context.Background(), email).Result()
//	if err != nil {
//		return "", fmt.Errorf("email code get cached err: %v", err)
//	}
//	return code, nil
//}
