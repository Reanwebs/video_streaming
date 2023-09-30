package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUniqueString() string {
	// Create a timestamp string using the current time
	timestamp := time.Now().UnixNano()

	// Generate a random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6 // Adjust the length as needed
	randomString := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	// Combine the timestamp and random string to create a unique identifier
	uniqueString := fmt.Sprintf("%d%s", timestamp, randomString)

	return uniqueString
}
