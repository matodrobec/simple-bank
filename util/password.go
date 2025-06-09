package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword( pass string) (string, error) {
    hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("faild to hash password")
    }
    return string(hashPass), nil
}

func CheckPassword( pass string, hashPass string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
}