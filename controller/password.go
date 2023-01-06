package controller

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(h)
}

func CompareHash(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err == nil {
		return true
	}
	return false
}
