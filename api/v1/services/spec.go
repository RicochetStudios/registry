package services

import (
	"fmt"
	"net/http"

	polarisv1alpha1 "github.com/RicochetStudios/polaris/apis/v1alpha1"
	"github.com/RicochetStudios/registry/api/v1/handlers"
	"github.com/gofiber/fiber/v2"
)

func CreateSpec() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Parse the request body.
		payload := new(polarisv1alpha1.ServerSpec)
		if err := ctx.BodyParser(payload); err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(
				handlers.ErrorResponse(fmt.Errorf("the request body was incorrect: %v", err)),
			)
		}

		// TODO: Convert the ServerSpec into a Agones GameServerSpec.

		return ctx.SendString("Create a new service")
	}
}
