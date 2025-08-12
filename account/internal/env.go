package internal

import (
	"log"
	"os"
	"strconv"
)

var (
	// The url of the main database
	EnvDbUrl string

	// String to be set on JWT token `iss` claim
	EnvTokenIssuer string
	// String to be used to sign JWT tokens
	EnvTokenSecret string
	// The lifetime of JWT tokens since its signing in seconds
	EnvTokenLifetime uint

	// The API endpoint to fetch data related to Indonesian provinces, regencies, cities, and more
	EnvIndoApiEndpoint string
)

func loadEnv() {
	EnvDbUrl = os.Getenv("DB_URL")
	EnvTokenIssuer = os.Getenv("TOKEN_ISSUER")
	EnvTokenSecret = os.Getenv("TOKEN_SECRET")

	tokenLifetime, err := strconv.ParseUint(os.Getenv("TOKEN_LIFETIME"), 10, 64)
	if err != nil {
		log.Fatalf("`TOKEN_LIFETIME`: %v", err)
	}
	EnvTokenLifetime = uint(tokenLifetime)

	EnvIndoApiEndpoint = os.Getenv("INDO_API_ENDPOINT")
}
