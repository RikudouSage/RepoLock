package service

import (
	"crypto/rand"
	"encoding/hex"
)

type RandomStringGenerator interface {
	GenerateRandomString(length int) string
}

func NewRandomStringGenerator() RandomStringGenerator {
	return &randomStringGenerator{}
}

type randomStringGenerator struct {
}

func (receiver *randomStringGenerator) GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	_, _ = rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
