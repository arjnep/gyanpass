package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func ComparePasswords(storedPassword string, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("unable to verify user password: %v", err)
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, fmt.Errorf("unable to verify user password: %v", err)
	}
	return hex.EncodeToString(shash) == pwsalt[0], nil
}
