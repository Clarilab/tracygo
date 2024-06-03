package atreugo

import (
	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

// CheckTracingIDs is a middleware for atreugo that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		correlationID := string(ctx.Request.Header.Peek(t.CorrelationIDKey()))
		requestID := string(ctx.Request.Header.Peek(t.RequestIDKey()))

		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// set userValues for resty middleware
		ctx.SetUserValue(t.CorrelationIDKey(), correlationID)
		ctx.SetUserValue(t.RequestIDKey(), requestID)

		ctx.Response.Header.Set(t.CorrelationIDKey(), correlationID)
		ctx.Response.Header.Set(t.RequestIDKey(), requestID)

		return ctx.Next()
	}
}
