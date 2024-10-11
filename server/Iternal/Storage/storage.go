package Storage

import (
	"context"
	"errors"
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
