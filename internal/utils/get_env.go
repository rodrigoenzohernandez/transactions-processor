package utils

import "os"

// getEnv retrieves the environment variable "ENV" and defaults to the provided fallback if not set.
func GetEnv(env string, fallback string) string {
	result := os.Getenv(env)
	if result == "" {
		result = fallback
	}
	return result
}
