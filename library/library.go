package library

var basePath string

func init() {
	// Get the environment variables.
	// TODO: Change default to "RicochetStudios/registry/main" before release.
	basePath = getEnvWithDefault("BASE_PATH", "RicochetStudios/registry/adopt-agones")
}

// getEnvWithDefault returns the value of the environment variable,
// or the default value if the environment variable is empty.
func getEnvWithDefault(env, def string) string {
	if env == "" {
		return def
	}
	return env
}
