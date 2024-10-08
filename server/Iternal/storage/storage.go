package storage

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
	"os"
	"server/Iternal/config"
)

type Storage struct {
	Db *pg.DB
}

func New(cfg config.DbConfig) (*Storage, error) {
	db := pg.Connect(&pg.Options{
		User:     cfg.Username,
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: cfg.DbName,
		Addr:     cfg.Address,
	})

	//TODO: init migrathions

	if err := db.Ping(context.Background()); err != nil {
		err = errors.New("failed connect to db: " + err.Error())
		return nil, err
	}
	return &Storage{Db: db}, nil
}
