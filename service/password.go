package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type PasswordVerifier interface {
	Verify(password, encodedHash string) (bool, error)
}

const (
	timeCost    = 3
	memoryCost  = 256 * 1024
	parallelism = 4

	saltLength = 16
	keyLength  = 32
)

type passwordManager struct {
}

func NewPasswordVerifier(hasher PasswordHasher) PasswordVerifier {
	return hasher.(*passwordManager)
}

func NewPasswordHasher() PasswordHasher {
	return &passwordManager{}
}

func (receiver *passwordManager) Verify(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
	}

	var memory, time uint32
	var parallel uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallel)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	test := argon2.IDKey([]byte(password), salt, time, memory, parallel, uint32(len(hash)))

	if len(test) != len(hash) {
		return false, nil
	}
	var diff byte
	for i := range test {
		diff |= test[i] ^ hash[i]
	}
	return diff == 0, nil
}

func (receiver *passwordManager) Hash(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memoryCost, timeCost, parallelism, b64Salt, b64Hash), nil
}
