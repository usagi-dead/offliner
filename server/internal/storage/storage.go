package storage

import (
	"context"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"server/internal/config"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var MigrationsFS embed.FS

type Storage struct {
	Db *pgxpool.Pool
}

func NewStorage(cfg config.DbConfig) (*Storage, error) {
	dbUrl := "postgres://" + cfg.Username + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + cfg.Address + "/" + cfg.DbName

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return &Storage{Db: pool}, nil
}

func (s *Storage) ApplyMigrations(cfg config.DbConfig) error {
	dbUrl := "postgres://" + cfg.Username + ":" + os.Getenv("POSTGRES_PASSWORD") + "@" + cfg.Address + "/" + cfg.DbName

	migrator := MustGetNewMigrator(MigrationsFS, migrationsDir)
	if err := migrator.ApplyMigrations(dbUrl); err != nil {
		return err
	}

	fmt.Println("Migrations applied!!")
	return nil
}
