package utils

import "golang.org/x/crypto/bcrypt"

func HashString(original string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(original), 8)
	return string(hash), err
}

func CompareHashAndOriginal(hash string, original string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(original)) == nil
}



