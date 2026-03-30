package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GeneratePasswordHash generates a bcrypt hash for a password (utility function)
// Run: go run hash_password.go to generate hashes
func GeneratePasswordHash() {
	password := "Reception@123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(hash))
}
