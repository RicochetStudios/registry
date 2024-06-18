package util

import (
	"os"
)

// RemapEnv sets game server environment variables, based on the given map of key-value pairs.
// It will get the value of an env var named after the key,
// and set an env var named the after pair value.
func RemapEnv(m map[string]string) error {
	for key, newKey := range m {
		if value, exists := os.LookupEnv(key); exists {
			InfoMessage("Found " + newKey + " , setting to " + value)
			if err := os.Setenv(newKey, value); err != nil {
				return err
			}
		}
	}
	return nil
}
