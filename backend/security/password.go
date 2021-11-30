package security

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"golang.org/x/crypto/scrypt"
)

//go:embed passwords_salt.txt
var salt string

func createHash(password string) ([]byte, error) {
	hash, err := scrypt.Key([]byte(password), []byte(salt), 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func CreatePasswordHash(password string) (string, error) {
	hash, err := createHash(password)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

func IsPassMatchHash(password, hash string) (bool, error) {
	storedHash, err := scrypt.Key([]byte(password), []byte(salt), 1<<15, 8, 1, 32) // todo: get params from hash
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}
	return bytes.Equal(storedHash, decodedHash), nil
}