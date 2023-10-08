package atreugo

import (
	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

// AtreugoCheckTracingIDs is a useBefore middleware for atreugo that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func AtreugoCheckTracingIDs(t *tracygo.TracyGo) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {

		correlationID := string(ctx.Request.Header.Peek(t.GetCorrelationID().String()))
		requestID := string(ctx.Request.Header.Peek(t.GetRequestID().String()))

		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// set userValue for resty middleware
		ctx.SetUserValue(t.GetCorrelationID(), correlationID)

		ctx.Response.Header.Set(t.GetCorrelationID().String(), correlationID)
		ctx.Response.Header.Set(t.GetRequestID().String(), requestID)

		return ctx.Next()
	}
}
