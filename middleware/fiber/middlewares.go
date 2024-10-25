package fiber

import (
	"context"

	"github.com/Clarilab/tracygo/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CheckTracingIDs is a middleware for fiber that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		correlationID := string(ctx.Request().Header.Peek(string(t.CorrelationIDKey())))
		requestID := string(ctx.Request().Header.Peek(string(t.RequestIDKey())))

		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Set values to UserContext
		userCtx := context.WithValue(ctx.UserContext(), t.CorrelationIDKey(), correlationID)
		userCtx = context.WithValue(userCtx, t.RequestIDKey(), requestID)
		ctx.SetUserContext(userCtx)

		// set userValues for resty middleware
		ctx.Context().SetUserValue(t.CorrelationIDKey(), correlationID)
		ctx.Context().SetUserValue(t.RequestIDKey(), requestID)

		ctx.Response().Header.Set(string(t.CorrelationIDKey()), correlationID)
		ctx.Response().Header.Set(string(t.RequestIDKey()), requestID)

		return ctx.Next()
	}
}
