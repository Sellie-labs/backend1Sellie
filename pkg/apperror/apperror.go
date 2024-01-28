package apperror

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(404, message)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(400, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(500, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(401, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(403, message)
}

func NewValidationError(message string) *AppError {
	return NewAppError(422, message)
}

// ErrorHandler is a middleware that converts AppError to fiber.Error.
func ErrorHandler(c *fiber.Ctx) error {
	// Try to execute the next middleware/handler
	err := c.Next()

	// Check if there was an error
	if err != nil {
		// Log the error, handle it, or send a custom response
		var fiberErr *fiber.Error
		if e, ok := err.(*AppError); ok {
			fiberErr = fiber.NewError(e.Code, e.Message)
		} else {
			fiberErr = fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		// Return the error to the client
		return c.Status(fiberErr.Code).SendString(fiberErr.Message)
	}

	// If no error, continue the chain
	return nil
}
