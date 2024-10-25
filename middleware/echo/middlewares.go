package echo

import (
	"context"

	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const tracyContextKey = "userCtx"

// CheckTracingIDs is a middleware for fiber that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(string(t.CorrelationIDKey()))
			requestID := c.Request().Header.Get(string(t.RequestIDKey()))

			if correlationID == "" {
				correlationID = uuid.NewString()
			}

			if requestID == "" {
				requestID = uuid.NewString()
			}

			tCtx := context.WithValue(context.Background(), t.CorrelationIDKey(), correlationID)
			tCtx = context.WithValue(tCtx, t.RequestIDKey(), requestID)
			c.Set(tracyContextKey, tCtx)

			c.Set(string(t.CorrelationIDKey()), correlationID)
			c.Set(string(t.RequestIDKey()), requestID)

			c.Response().Header().Set(string(t.CorrelationIDKey()), correlationID)
			c.Response().Header().Set(string(t.RequestIDKey()), requestID)

			return next(c)
		}
	}
}

// GetUserContext is a helper function to extract a context set by tracygo from a echo.Context
// This mirrors context attachment functionality of other libs
func GetUserContext(c echo.Context) context.Context {
	if val := c.Get(tracyContextKey); val != nil {
		if ctx, ok := val.(context.Context); ok {
			return ctx
		}
	}

	return nil
}
