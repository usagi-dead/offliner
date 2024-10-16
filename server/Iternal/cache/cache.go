package cache

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"server/Iternal/config"
	"strings"
	"time"
)

type Cache struct {
	db                           *redis.Client
	stateExpiration              time.Duration
	emailConfirmedCodeExpiration time.Duration
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

func (c Cache) CreateEmailConfirmedCode(email string) (string, error) {
	code, err := generateCode(6)
	if err != nil {
		return "", fmt.Errorf("generate code err: %v", err)
	}
	if err := c.db.Set(context.Background(), email, code, c.emailConfirmedCodeExpiration).Err(); err != nil {
		return "", fmt.Errorf("email code set cached err: %v", err)
	}
	return code, nil
}

func (c Cache) GetEmailConfirmedCode(email string) (string, error) {
	code, err := c.db.Get(context.Background(), email).Result()
	if err != nil {
		return "", fmt.Errorf("email code get cached err: %v", err)
	}
	return code, nil
}

func generateCode(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)[:length]), nil
}
