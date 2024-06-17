package util

import (
	"fmt"
	"os"
)

// GetEnvWithDefault returns the value of the environment variable with the given key.
// If the environment variable is not set, it returns the default value.
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// RemapEnv sets game server environment variables, based on the given map of key-value pairs.
// It will get the value of an env var named after the key,
// and set an env var named the after pair value.
func RemapEnv(m map[string]string) error {
	for key, newKey := range m {
		if value, exists := os.LookupEnv(key); exists {
			fmt.Println("Setting", newKey, "to", value)
			if err := os.Setenv(newKey, value); err != nil {
				return err
			}
		} else {
			fmt.Println("Key", key, "not found")
		}
	}
	return nil
}
