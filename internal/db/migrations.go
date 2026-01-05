package db

import (
	"github.com/levinOo/go-crudl-task/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Выполнение миграций БД
func RunMigrations(pg *Postgres) error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	poolConfig := pg.Pool.Config()

	stdDB := stdlib.OpenDB(*poolConfig.ConnConfig)
	defer stdDB.Close()

	if err := goose.Up(stdDB, "."); err != nil {
		return err
	}

	return nil
}
