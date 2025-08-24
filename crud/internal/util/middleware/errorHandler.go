package middleware

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops/adapter"
	"github.com/solsteace/goody/lib/payload"
)

// https://docs.gofiber.io/guide/error-handling/
type errorHandler struct {
	payloader payload.Loader
}

func NewErrorHandler(payloader payload.Loader) errorHandler {
	return errorHandler{payloader: payloader}
}

func (eh errorHandler) Handle(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	statusCode := http.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) { // fiber specific error
		statusCode = e.Code
	} else {
		statusCode = adapter.HttpStatusCode(err)
	}

	return c.
		Status(statusCode).
		JSON(eh.payloader.Err(c.Method(), []error{err}))
}
