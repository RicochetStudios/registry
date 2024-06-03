package handlers

import "github.com/gofiber/fiber/v2"

// ErrorResponse is the singular ErrorResponse that will be passed in the response by handlers.
func ErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}
