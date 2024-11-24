package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"server/internal/cache"
	u "server/internal/features/user"
)

type UserQuery struct {
	log *slog.Logger
	db  *pgxpool.Pool
	ch  *cache.Cache
}

func NewUserQuery(log *slog.Logger, db *pgxpool.Pool, ch *cache.Cache) u.UserData {
	return &UserQuery{
		log: log,
		db:  db,
		ch:  ch,
	}
}

func (uq *UserQuery) CreateUser(email string, hashedPassword string) (int64, error) {
	var userID int64
	err := uq.db.QueryRow(context.Background(),
		`INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING user_id`, email, hashedPassword).Scan(&userID)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return 0, u.ErrEmailExists
		}
		return 0, err
	}

	return userID, nil
}

func (uq *UserQuery) CreateOauthUser(user *u.User) (int64, error) {
	query := `
        INSERT INTO users (
            email, hashed_password, role, surname, name, patronymic, 
            date_of_birth, phone_number, avatar_url, verified_email, gender
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        ) RETURNING user_id`

	var userID int64
	err := uq.db.QueryRow(
		context.Background(),
		query,
		user.Email,
		nil,
		user.Role,
		user.Surname,
		user.Name,
		user.Patronymic,
		user.DateOfBirth,
		user.PhoneNumber,
		user.AvatarUrl,
		user.VerifiedEmail,
		user.Gender,
	).Scan(&userID)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return 0, u.ErrEmailExists
		}
		return 0, err
	}
	return userID, nil
}

func (uq *UserQuery) GetUserByEmail(Email string) (*u.User, error) {
	row := uq.db.QueryRow(context.Background(),
		`SELECT * FROM users WHERE email = $1`,
		Email,
	)

	user := &u.User{}

	err := row.Scan(
		&user.UserId,
		&user.HashedPassword,
		&user.Role,
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.DateOfBirth,
		&user.PhoneNumber,
		&user.Email,
		&user.AvatarUrl,
		&user.VerifiedEmail,
		&user.Gender,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, u.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (uq *UserQuery) SaveStateCode(state string) error {
	if err := uq.ch.Db.Set(context.Background(), state, "true", uq.ch.StateExpiration).Err(); err != nil {
		return fmt.Errorf("state token set cached err: %v", err)
	}
	return nil
}

func (uq *UserQuery) VerifyStateCode(state string) (bool, error) {
	state, err := uq.ch.Db.Get(context.Background(), state).Result()
	if err != nil {
		return false, fmt.Errorf("state token get cached err: %v", err)
	}

	if state == "true" {
		if err := uq.ch.Db.Del(context.Background(), state).Err(); err != nil {
			return false, fmt.Errorf("failed to delete state token from cache: %v", err)
		}
		return true, nil
	}

	return false, nil
}
