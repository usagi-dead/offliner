package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"server/internal/config"
	"server/internal/storage/models"
)

var (
	ErrEmailExists    = errors.New("users_email_key")
	ErrEmailNotExists = errors.New("user_email_not_exists")
)

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg config.DbConfig) (*Storage, error) {
	dbUrl := "postgres://" + cfg.Username + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + cfg.Address + "/" + cfg.DbName

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateUser(email string, hashedPassword string) error {
	_, err := s.db.Exec(context.Background(),
		`INSERT INTO users (email, hashed_password) VALUES ($1, $2)`, email, hashedPassword)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return ErrEmailExists
		}
		return err
	}
	return nil
}

func (s *Storage) GetUserByEmail(Email string) (*models.User, error) {

	row := s.db.QueryRow(context.Background(),
		`SELECT * FROM users WHERE email = $1`,
		Email,
	)

	user := &models.User{}

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

	// Обрабатываем ошибки
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", Email)
		}
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return user, nil
}

func (s *Storage) GetUserById(UserId int64) (*models.User, error) {
	row := s.db.QueryRow(context.Background(),
		`SELECT * FROM users WHERE user_id = $1`,
		UserId,
	)

	user := &models.User{}

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
		&user.Gender,
	)

	// Обрабатываем ошибки
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %v not found", UserId)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (s *Storage) IsEmailConfirmed(UserId int64) (bool, error) {
	row := s.db.QueryRow(context.Background(),
		`SELECT verified_email FROM users WHERE user_id = $1`,
		UserId,
	)

	var EmailStatus string
	err := row.Scan(&EmailStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("user with id %v not found", UserId)
		}
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	return EmailStatus == "confirmed", nil
}

func (s *Storage) UpdateEmailStatus(Email string) error {
	_, err := s.db.Exec(context.Background(),
		`UPDATE users SET verified_email = true WHERE email = $1`,
		Email,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return ErrEmailNotExists
		}
		return err
	}

	return nil
}
