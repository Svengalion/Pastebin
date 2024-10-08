package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	HashSize     int    = 8
	HashAlphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateHash() (hash string, err error) {
	return gonanoid.Generate(HashAlphabet, HashSize)
}
