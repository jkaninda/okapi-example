package utils

import "os"

func GetSingingSecret() string {
	value := os.Getenv("JWT_SIGNING_SECRET")
	if value == "" {
		return "supersecret"
	}
	return value
}
