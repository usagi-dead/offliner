package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"server/internal/config"
	"server/internal/lib/CustomErrors"
	"server/internal/storage/models"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg config.DbConfig) (*Storage, error) {
	dbUrl := "postgres://" + cfg.Username + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + cfg.Address + "/" + cfg.DbName

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	migrator := MustGetNewMigrator(MigrationsFS, migrationsDir)
	if err := migrator.ApplyMigrations(dbUrl); err != nil {
		panic(err)
	}

	fmt.Printf("Migrations applied!!")

	return &Storage{db: pool}, nil
}

func (s *Storage) CreateUser(email string, hashedPassword string) (int64, error) {
	var userID int64
	err := s.db.QueryRow(context.Background(),
		`INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING user_id`, email, hashedPassword).Scan(&userID)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return 0, CustomErrors.ErrEmailExists
		}
		return 0, err
	}
	return userID, nil
}

func (s *Storage) CreateUserAfterOauth(user *models.User) (int64, error) {
	query := `
        INSERT INTO users (
            email, hashed_password, role, surname, name, patronymic, 
            date_of_birth, phone_number, avatar_url, verified_email, gender
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        ) RETURNING user_id`

	var userID int64
	err := s.db.QueryRow(
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
			return 0, CustomErrors.ErrEmailExists
		}
		return 0, err
	}
	return userID, nil
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
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return user, nil
}

func (s *Storage) UpdateUserById(user *models.User) error {
	var exists bool
	err := s.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", user.UserId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return CustomErrors.ErrEmailNotExists
	}

	query := `
        UPDATE users 
        SET 
            hashed_password = COALESCE($1, hashed_password),
            role = $2,
            surname = COALESCE($3, surname),
            name = COALESCE($4, name),
            patronymic = COALESCE($5, patronymic),
            date_of_birth = COALESCE($6, date_of_birth),
            phone_number = COALESCE($7, phone_number),
            avatar_url = COALESCE($8, avatar_url),
            gender = COALESCE($9, gender)
        WHERE user_id = $10`

	_, err = s.db.Exec(context.Background(),
		query,
		user.HashedPassword,
		user.Role,
		user.Surname,
		user.Name,
		user.Patronymic,
		user.DateOfBirth,
		user.PhoneNumber,
		user.AvatarUrl,
		user.Gender,
		user.UserId,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *Storage) GetUserById(UserId int64) error {
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
			return fmt.Errorf("user with id %v not found", UserId)
		}
		return fmt.Errorf("failed to get user: %v", err)
	}

	return nil
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
			return CustomErrors.ErrEmailNotExists
		}
		return err
	}

	return nil
}
