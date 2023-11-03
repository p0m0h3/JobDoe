package handlers

import (
	"errors"

	"git.fuzz.codes/fuzzercloud/workerengine/schemas"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(schemas.ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	})
}

func UnauthorizedError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(schemas.ErrorResponse{
		Code:    fiber.ErrUnauthorized.Code,
		Message: err.Error(),
	})
}

func InternalServerError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(schemas.ErrorResponse{
		Code:    fiber.ErrInternalServerError.Code,
		Message: fiber.ErrInternalServerError.Message,
	})
}

func ForbiddenError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(schemas.ErrorResponse{
		Code:    fiber.ErrForbidden.Code,
		Message: fiber.ErrForbidden.Message,
	})
}

func NotFoundError(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(schemas.ErrorResponse{
		Code:    fiber.ErrNotFound.Code,
		Message: fiber.ErrNotFound.Message,
	})
}

func BadRequestError(c *fiber.Ctx, errors []error) error {
	return c.Status(fiber.StatusBadRequest).JSON(schemas.ErrorResponse{
		Code:    fiber.ErrBadRequest.Code,
		Message: fiber.ErrBadRequest.Message,
		Errors:  errors,
	})
}
