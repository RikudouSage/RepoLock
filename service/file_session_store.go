package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"go.chrastecky.dev/repolock/config/data"
)

type FileSessionStore interface {
	scs.Store
}

func NewFileSessionStore(
	config *data.GlobalConfig,
) FileSessionStore {
	return &fileSessionStore{
		sessionStorePath: config.SessionStorePath,
		now:              time.Now,
	}
}

type fileSessionStore struct {
	now              func() time.Time
	sessionStorePath string
}

func (receiver *fileSessionStore) Delete(token string) (err error) {
	path := receiver.getSessionPath(token)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}

func (receiver *fileSessionStore) Find(token string) (bytes []byte, found bool, err error) {
	path := receiver.getSessionPath(token)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}

		return nil, false, err
	}
	defer file.Close()

	lengthBuf := make([]byte, 2)
	if _, err := io.ReadFull(file, lengthBuf); err != nil {
		return nil, false, fmt.Errorf("failed reading expiry length: %w", err)
	}

	expiryLength, err := strconv.Atoi(string(lengthBuf))
	if err != nil {
		return nil, false, fmt.Errorf("failed parsing expiry length: %w", err)
	}

	expiryBuf := make([]byte, expiryLength)
	if _, err := io.ReadFull(file, expiryBuf); err != nil {
		return nil, false, fmt.Errorf("failed reading expiry: %w", err)
	}

	expiry, err := time.Parse(time.RFC3339, string(expiryBuf))
	if err != nil {
		return nil, false, fmt.Errorf("failed parsing expiry: %w", err)
	}

	if expiry.Before(receiver.now()) {
		defer receiver.Delete(token)
		return nil, false, nil
	}

	result, err := io.ReadAll(file)
	if err != nil {
		return nil, false, fmt.Errorf("failed reading session data: %w", err)
	}

	return result, true, nil
}

func (receiver *fileSessionStore) Commit(token string, bytes []byte, expiry time.Time) (err error) {
	path := receiver.getSessionPath(token)
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return fmt.Errorf("failed creating session directory: %w", err)
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed creating session file: %w", err)
	}
	defer file.Close()

	expiryStr := expiry.Format(time.RFC3339)

	if _, err = file.Write([]byte(strconv.Itoa(len(expiryStr)))); err != nil {
		return fmt.Errorf("failed writing data: %w", err)
	}
	if _, err = file.Write([]byte(expiryStr)); err != nil {
		return fmt.Errorf("failed writing data: %w", err)
	}
	if _, err = file.Write(bytes); err != nil {
		return fmt.Errorf("failed writing data: %w", err)
	}

	return nil
}

func (receiver *fileSessionStore) getSessionPath(token string) string {
	return filepath.Join(receiver.sessionStorePath, token+".session")
}
