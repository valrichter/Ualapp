package util

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

// Generates a random full name
func RandomFullName() string {
	fullName := gofakeit.FirstName() + " " + gofakeit.LastName()
	return fullName
}

// Generates a random amount of money
func RandomMoney() int64 {
	return int64(gofakeit.IntRange(0, 1000))
}

// Generates a random email
func RandomEmail() string {
	email := gofakeit.Email()
	return email
}

// RandomInt generates a random positive integer between min and max
func RandomInt(min, max int) int {
	n := rand.Intn(max-min+1) + min
	return n
}

// RandomPassword generates a random password
func RandomPassword(length int) string {
	password := gofakeit.Password(false, false, false, false, false, length)
	return password
}

func RandomUsername() string {
	username := gofakeit.Username()
	return username
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomString generates a random string of length n
func RandomString(n int) string {
	bits := []rune{}
	k := len(alphabet)

	for i := 0; i < n; i++ {
		index := rand.Intn(k)
		bits = append(bits, rune(alphabet[index]))
	}
	return string(bits)
}
