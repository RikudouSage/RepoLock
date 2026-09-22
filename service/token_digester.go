package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"go.chrastecky.dev/repolock/config/data"
)

type TokenDigester interface {
	Digest(token string) string
}

func NewTokenDigester(config *data.GlobalConfig) TokenDigester {
	return &tokenDigester{key: []byte(config.Secret)}
}

type tokenDigester struct {
	key []byte
}

func (receiver *tokenDigester) Digest(token string) string {
	mac := hmac.New(sha256.New, receiver.key)
	_, _ = mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
