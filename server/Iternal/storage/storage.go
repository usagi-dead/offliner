package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"server/Iternal/config"
	"server/Iternal/storage/models"
)

var (
	ErrEmailExists = errors.New("users_email_key")
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

	//TODO: init migrathions

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateUser(user *models.User) error {
	//err := s.db.QueryRow(context.Background(),
	//	`INSERT INTO `)

	err := s.db.QueryRow(context.Background(),
		`INSERT INTO users (hashed_password, surname, name, patronymic, date_of_birth, phone_number, email, gender, role)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
        RETURNING user_id`,
		user.HashedPassword, user.Surname, user.Name, user.Patronymic, user.DateOfBirth, user.PhoneNumber, user.Email, user.Gender, user.Role).Scan(&user.UserId)

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
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return user, nil
}
