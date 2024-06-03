package fiber

import (
	"github.com/Clarilab/tracygo/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CheckTracingIDs is a middleware for fiber that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		correlationID := string(ctx.Request().Header.Peek(t.CorrelationIDKey()))
		requestID := string(ctx.Request().Header.Peek(t.RequestIDKey()))

		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// set userValues for resty middleware
		ctx.Context().SetUserValue(t.CorrelationIDKey(), correlationID)
		ctx.Context().SetUserValue(t.RequestIDKey(), requestID)

		ctx.Response().Header.Set(t.CorrelationIDKey(), correlationID)
		ctx.Response().Header.Set(t.RequestIDKey(), requestID)

		return ctx.Next()
	}
}
