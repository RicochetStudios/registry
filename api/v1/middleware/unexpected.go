package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/RicochetStudios/registry/api/v1/handlers"
)

// HitTest checks if the api server is running.
func UnexpectedRouter() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Status(http.StatusNotFound)
		return ctx.JSON(handlers.UnexpectedErrorResponse())
	}
}
