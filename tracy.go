// The TracyGo package provides functionality for tracing
// a correlation identifier through multiple go microservices.
package tracygo

import (
	"context"
)

const (
	correlationID = "X-Correlation-ID"
	requestID     = "X-Request-ID"
)

// TracyGo is a struct for the tracy object.
type TracyGo struct {
	correlationID CorrelationID
	requestID     RequestID
}

// New creates a new TracyGo object and uses the options on it.
func New(options ...Option) *TracyGo {
	tracy := &TracyGo{
		correlationID: correlationID,
		requestID:     requestID,
	}

	for _, option := range options {
		option(tracy)
	}

	return tracy
}

// Get retrieves the underlying correlationID key.
func (t *TracyGo) GetCorrelationID() CorrelationID {
	return t.correlationID
}

// GetRequestID retrieves the underlying requestID key.
func (t *TracyGo) GetRequestID() RequestID {
	return t.requestID
}

// FromContext returns the correlationID from the given context, or the an empty string.
func (t *TracyGo) FromContext(ctx context.Context) string {
	if ctx != nil {
		if correlationID, ok := ctx.Value(t.correlationID).(string); ok {
			return correlationID
		}
	}

	return ""
}

// NewContext sets the correlationID to use in the given context. If ctx is nil, a new context is created.
func (t *TracyGo) NewContext(ctx context.Context, correlationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, t.correlationID, correlationID)
}
