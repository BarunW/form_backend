package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log/slog"
	"strings"
)

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		slog.Error("Failed to generate crypto rand", "details", err.Error())
		return nil, err
	}
	return salt, nil
}

func hashWithSalt(data []byte, salt []byte) string {
	combinedData := append(data, salt...)
	hasher := sha256.New()
	hasher.Write(combinedData)
	return strings.ToLower(base32.StdEncoding.EncodeToString(hasher.Sum(nil))[:11])
}

func GenerateAccountId() (string, error) {
	inputData, err := GenerateULID()
	if err != nil {
		slog.Error("Failed to generate ulid", "details", err.Error())
		return "", err
	}

	salt, err := generateSalt(16) // 16-byte salt
	if err != nil {
		return "", err
	}
	// return the hasehed userId
	// encoded with bas32
	return hashWithSalt([]byte(inputData), salt), nil

}
