package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HandleError(message string, err error) {
	if err != nil {
		log.Fatal(message, err.Error())
	}
}

func HashPassword(p string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	HandleError("password hashing error", err)
	return string(h)
}

func CompareHash(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}
