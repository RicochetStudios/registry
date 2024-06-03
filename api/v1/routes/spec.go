package routes

import (
	"github.com/RicochetStudios/registry/api/v1/services"
	"github.com/gofiber/fiber/v2"
)

func SpecRouter(app fiber.Router) {
	// Create the schema group.
	app = app.Group("/spec")

	// Define the routes.
	app.Post("/create", services.CreateSpec())
}
