package data

import (
	"fmt"
	"os"
	"strings"
)

type DatabaseType string

func (receiver DatabaseType) IsValid() bool {
	return receiver == DatabaseTypeSQLite || receiver == DatabaseTypePostgres
}

const (
	DatabaseTypeSQLite   DatabaseType = "sqlite"
	DatabaseTypePostgres DatabaseType = "postgres"
)

type GlobalConfig struct {
	// runtime
	Port     uint16 `default:"8080"`
	Timezone string `default:"UTC"`
	Debug    bool   `default:"false"`

	// config
	Secret               string `required:"true" split_words:"true"`
	RegistrationsEnabled bool   `default:"true" split_words:"true"`
	PasswordLoginEnabled bool   `default:"true" split_words:"true"`
	AnyoneCanRequestJoin bool   `default:"true" split_words:"true"`
	SessionStorePath     string `default:"$HOME/.config/repo-lock/sessions" split_words:"true"`

	// db
	DatabaseType     DatabaseType `default:"sqlite" split_words:"true"`
	DatabasePath     string       `default:"./data.sqlite3" split_words:"true"`
	DatabaseHost     string       `split_words:"true"`
	DatabaseName     string       `split_words:"true"`
	DatabaseUser     string       `split_words:"true"`
	DatabasePassword string       `split_words:"true"`
	DatabasePort     uint16       `split_words:"true" default:"5432"`
	DatabaseSSLMode  string       `default:"require" split_words:"true"`
}

func (receiver *GlobalConfig) Validate() error {
	if !receiver.DatabaseType.IsValid() {
		return fmt.Errorf("database type '%s' is not supported", receiver.DatabaseType)
	}

	if receiver.DatabaseType == DatabaseTypePostgres && (receiver.DatabaseHost == "" || receiver.DatabaseName == "") {
		return fmt.Errorf("database host or name is empty for postgres mode")
	}

	if len(receiver.Secret) < 32 {
		return fmt.Errorf("app secret must contain at least 32 characters (set APP_SECRET)")
	}

	return nil
}

func (receiver *GlobalConfig) Normalize() error {
	if strings.Contains(receiver.SessionStorePath, "$HOME") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("the session store path contains $HOME, but we could not resolve the user's home dir: %w", err)
		}
		receiver.SessionStorePath = strings.ReplaceAll(receiver.SessionStorePath, "$HOME", home)
	}

	return nil
}
