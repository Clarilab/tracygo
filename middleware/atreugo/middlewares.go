package atreugo

import (
	"context"

	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

// CheckTracingIDs is a middleware for atreugo that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(tracy *tracygo.TracyGo) func(request *atreugo.RequestCtx) error {
	return func(request *atreugo.RequestCtx) error {
		correlationID := string(request.Request.Header.Peek(string(tracy.CorrelationIDKey())))
		requestID := string(request.Request.Header.Peek(string(tracy.RequestIDKey())))

		if correlationID == "" {
			correlationID = uuid.NewString()
		}

		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Set values to attachedContext. While the request fulfills the context interface, it is not recommended for performance and opens up some pitfalls.
		aCtx := request.AttachedContext()
		if aCtx == nil {
			aCtx = context.Background()
		}

		aCtx = context.WithValue(aCtx, tracy.CorrelationIDKey(), correlationID)
		aCtx = context.WithValue(aCtx, tracy.RequestIDKey(), requestID)

		request.AttachContext(aCtx)

		// set userValues for resty middleware (legacy)
		request.SetUserValue(tracy.CorrelationIDKey(), correlationID)
		request.SetUserValue(tracy.RequestIDKey(), requestID)

		request.Response.Header.Set(string(tracy.CorrelationIDKey()), correlationID)
		request.Response.Header.Set(string(tracy.RequestIDKey()), requestID)

		return request.Next()
	}
}
