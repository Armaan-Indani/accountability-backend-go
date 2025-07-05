package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

// Config func to get env value
func Config(key string) string {
	// load .env file only once and ignore errors in production
	once.Do(func() {
		// Try to load .env file, but don't panic if it doesn't exist
		// This is useful for local development
		godotenv.Load(".env")
	})
	return os.Getenv(key)
}
