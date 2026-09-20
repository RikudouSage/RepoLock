package config

import (
	"database/sql"
	"embed"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/pressly/goose/v3"
	"go.chrastecky.dev/repolock/config/data"
	"go.chrastecky.dev/repolock/migrations/postgres"
	"go.chrastecky.dev/repolock/migrations/sqlite"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

type gooseLogger struct {
	*zap.Logger
}

func (receiver *gooseLogger) Fatalf(format string, v ...any) {
	receiver.Fatal(fmt.Sprintf(format, v...))
}

func (receiver *gooseLogger) Printf(format string, v ...any) {
	receiver.Info(fmt.Sprintf(format, v...))
}

func createPostgresDsn(config *data.GlobalConfig) string {
	uri := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(config.DatabaseHost, strconv.Itoa(int(config.DatabasePort))),
		Path:   config.DatabaseName,
	}

	if config.DatabaseUser != "" {
		if config.DatabasePassword != "" {
			uri.User = url.UserPassword(config.DatabaseUser, config.DatabasePassword)
		} else {
			uri.User = url.User(config.DatabaseUser)
		}
	}

	query := uri.Query()
	query.Set("sslmode", config.DatabaseSSLMode)
	query.Set("application_name", "repo_lock")
	query.Set("connect_timeout", "5")
	uri.RawQuery = query.Encode()

	return uri.String()
}

func migrate(db *sql.DB, dbType data.DatabaseType, migrations embed.FS, logger goose.Logger) error {
	goose.SetLogger(logger)
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect(string(dbType)); err != nil {
		return fmt.Errorf("invalid database type: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("could not migrate db: %w", err)
	}

	return nil
}

func createDatabase(
	config *data.GlobalConfig,
	logger *zap.Logger,
) (db *sql.DB, err error) {
	var migrations embed.FS

	if config.DatabaseType == data.DatabaseTypeSQLite {
		migrations = sqlite.Migrations
		db, err = sql.Open(
			"sqlite3",
			fmt.Sprintf(
				"file:%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL",
				config.DatabasePath,
			),
		)

		if err != nil {
			err = fmt.Errorf("failed to open sqlite database: %w", err)
		}
	} else {
		migrations = postgres.Migrations
		db, err = sql.Open("postgres", createPostgresDsn(config))
		if err != nil {
			err = fmt.Errorf("failed to open postgres database: %w", err)
		}
	}

	if err != nil {
		return
	}

	err = migrate(db, config.DatabaseType, migrations, &gooseLogger{logger})

	return
}

func provideDatabase() fx.Option {
	return fx.Module(
		"database",
		fx.Provide(createDatabase),
	)
}
