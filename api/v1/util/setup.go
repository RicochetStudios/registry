package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RicochetStudios/registry/api/v1/middleware"
	"github.com/RicochetStudios/registry/api/v1/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/joho/godotenv"
)

const ApiGroup = "api/v1"

// SetupAPI initializes and starts the API.
func SetupAPI() *fiber.App {
	// Load the .env file.
	err := getDotEnv()
	if err != nil {
		panic(err)
	}

	// Create the API.
	app := fiber.New()
	app.Use(cors.New())

	// Set the API version group.
	api := app.Group(ApiGroup)

	// Run the routers.
	routes.SpecRouter(api)

	// Middleware to handle non-existent endpoints
	app.Use(middleware.UnexpectedRouter())

	return app
}

// getDotEnv reads a dotenv file from the root directory
// and loads it's contents into the environment.
func getDotEnv() error {
	// Get the root directory of the executable.
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error finding .env root directory:", err)
	}
	root := filepath.Dir(ex)

	// Set the path to the .env file.
	dotEnvConfigPath := filepath.Join(root, ".env")

	// Load the .env file.
	err = godotenv.Load(dotEnvConfigPath)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	return nil
}

// GetEnv returns the value of an environment variable.
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

// GetEnvWithDefault returns the value of an environment variable,
// or a default value if the environment variable is not set.
func GetEnvWithDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
