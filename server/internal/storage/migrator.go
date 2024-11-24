package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // Важно для поддержки stdlib
)

type Migrator struct {
	srcDriver source.Driver
}

func MustGetNewMigrator(sqlFiles embed.FS, dirName string) *Migrator {
	srcDriver, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		panic(fmt.Errorf("failed to initialize source driver: %w", err))
	}
	return &Migrator{srcDriver: srcDriver}
}

func (m *Migrator) ApplyMigrations(dbURL string) error {
	// Открываем соединение через stdlib
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("unable to open database: %w", err)
	}
	defer db.Close()

	// Создаем драйвер миграций для pgx
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("unable to create pgx driver: %w", err)
	}

	// Создаем мигратор
	migrator, err := migrate.NewWithInstance(
		"iofs", m.srcDriver, "pgx", driver,
	)
	if err != nil {
		return fmt.Errorf("unable to create migrator: %w", err)
	}
	defer migrator.Close()

	// Применяем миграции
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
