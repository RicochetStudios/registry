package util

import "os"

// GetEnvWithDefault returns the value of the environment variable with the given key.
// If the environment variable is not set, it returns the default value.
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// SetEnv sets the environment variable with the given key to the given value.
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// RemapEnv sets game server environment variables, based on the given map of key-value pairs.
// It will get the value of an env var named after the key,
// and set an env var named the after pair value.
func RemapEnv(m map[string]string) error {
	for key, value := range m {
		if envValue, exists := os.LookupEnv(key); exists {
			if err := os.Setenv(value, envValue); err != nil {
				return err
			}
		}
	}
	return nil
}
