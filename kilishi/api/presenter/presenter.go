package presenter

import "github.com/gofiber/fiber/v2"

func SuccessResponse(message string, data any) *fiber.Map {
	return &fiber.Map{
		"message": message,
		"data":    data,
	}
}

func ErrorResponse(message string, errors ...error) *fiber.Map {
	return &fiber.Map{
		"message": message,
		"errors":  errors,
	}
}
