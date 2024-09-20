package echo

import (
	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CheckTracingIDs is a middleware for fiber that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(tracer *tracygo.TracyGo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(tracer.CorrelationIDKey())
			requestID := c.Request().Header.Get(tracer.RequestIDKey())

			if correlationID == "" {
				correlationID = uuid.NewString()
			}

			if requestID == "" {
				requestID = uuid.NewString()
			}

			c.Set(tracer.CorrelationIDKey(), correlationID)
			c.Set(tracer.RequestIDKey(), requestID)

			c.Response().Header().Set(tracer.CorrelationIDKey(), correlationID)
			c.Response().Header().Set(tracer.RequestIDKey(), requestID)

			return next(c)
		}
	}
}
