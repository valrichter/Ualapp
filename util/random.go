package util

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

// Generates a random full name
func RandomFullName() string {
	fullName := gofakeit.FirstName() + " " + gofakeit.LastName()
	return fullName
}

// Generates a random amount of money
func RandomMoney(min, max int) int64 {
	return int64(gofakeit.IntRange(min, max))
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
	password := gofakeit.Password(true, true, true, true, true, length)
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

func GenerateAccountNumber(accountID int32, currency string) (string, error) {
	config, ok := Currencies[currency]
	if !ok {
		return "", fmt.Errorf("invalid currency: %s", currency)
	}
	activeTime := time.Now().Format("20060102150405")
	initialValue := fmt.Sprintf("%s%d", config.Id, accountID)
	// account number should be 10 in length
	finalValue := ""
	reminder := 10 - len(initialValue)
	if reminder > 0 {
		finalValue = activeTime[:reminder]
	}
	accountNumber := fmt.Sprintf("%s%s", initialValue, finalValue)
	return accountNumber, nil
}
