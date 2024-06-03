package handlers

import "github.com/gofiber/fiber/v2"

// SetupErrorResponse is the singular ErrorResponse that will be passed in the response by handlers.
func UnexpectedErrorResponse() *fiber.Map {
	return &fiber.Map{
		"status":  404,
		"message": "Wake up Mr Freeman.",
		"error":   "This is not an endpoint.",
	}
}
