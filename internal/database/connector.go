package database

import (
	"crypto/tls"
	"database/sql"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
)

func Connect(c *config.ConnConfig) (*sql.DB, error) {
	pgxConfig, err := pgx.ParseEnvLibpq()
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if c.SslMode == "require" {
		tlsConfig = &tls.Config{}
	}
	optConfig := pgx.ConnConfig{
		Host:      c.Host,
		Port:      uint16(c.Port),
		Database:  c.DbName,
		User:      c.User,
		Password:  c.Password,
		TLSConfig: tlsConfig,
	}

	pgxConfig = pgxConfig.Merge(optConfig)

	return stdlib.OpenDB(pgxConfig), nil
}
