package sql

import (
	"embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func ApplyMigrations(pool *pgxpool.Pool) {
	// Set up Goose
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to set dialect", "error", err)
		return
	}

	// Run migrations - use empty string as path since we're using embed.FS
	if err := goose.Up(stdlib.OpenDBFromPool(pool), "migrations"); err != nil {
		slog.Error("Migration failed", "error", err)
	}
}
