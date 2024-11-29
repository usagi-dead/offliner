package data

import (
	"context"
	"server/internal/cache"
	u "server/internal/features/user"
)

type UserCacheClient interface {
	SaveStateCode(state string) error
	VerifyStateCode(state string) (bool, error)
	SaveEmailConfirmedCode(email string, code string) error
	GetEmailConfirmedCode(email string) (string, error)
}

type UserCache struct {
	ch *cache.Cache
}

func NewUserCache(ch *cache.Cache) *UserCache {
	return &UserCache{
		ch: ch,
	}
}

func (c *UserCache) SaveStateCode(state string) error {
	if err := c.ch.Db.Set(context.Background(), state, "true", c.ch.StateExpiration).Err(); err != nil {
		return u.ErrInternal
	}
	return nil
}

func (c *UserCache) VerifyStateCode(state string) (bool, error) {
	state, err := c.ch.Db.Get(context.Background(), state).Result()
	if err != nil {
		return false, u.ErrInvalidState
	}

	if state == "true" {
		if err := c.ch.Db.Del(context.Background(), state).Err(); err != nil {
			return false, u.ErrInternal
		}
		return true, nil
	}

	return false, nil
}

func (c *UserCache) SaveEmailConfirmedCode(email string, code string) error {
	if err := c.ch.Db.Set(context.Background(), email, code, c.ch.EmailConfirmedCodeExpiration).Err(); err != nil {
		return u.ErrInternal
	}
	return nil
}

func (c *UserCache) GetEmailConfirmedCode(email string) (string, error) {
	code, err := c.ch.Db.Get(context.Background(), email).Result()
	if err != nil {
		return "", u.ErrInternal
	}
	return code, nil
}
