package Storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"server/Iternal/Storage/models"
	"server/Iternal/config"
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

func (s *Storage) CreateUser(user models.User) error {
	_, err := s.db.Exec(context.Background(),
		`INSERT INTO users (hashed_password, surname, name, patronymic, date_of_birth, phone_number, email, gender, role)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.HashedPassword, user.Surname, user.Name, user.Patronymic, user.DateOfBirth, user.PhoneNumber, user.Email, user.Gender, user.Role)

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
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", Email)
		}
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return user, nil
}
