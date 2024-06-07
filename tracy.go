// The TracyGo package provides functionality for tracing
// a correlation identifier through multiple go microservices.
package tracygo

import (
	"context"

	"github.com/google/uuid"
)

const (
	correlationID = "X-Correlation-ID"
	requestID     = "X-Request-ID"
)

// TracyGo is a struct for the tracy object.
type TracyGo struct {
	correlationID string
	requestID     string
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

// CorrelationIDKey returns the underlying correlation id key.
func (t *TracyGo) CorrelationIDKey() string {
	return t.correlationID
}

// RequestIDKey returns the underlying request id key.
func (t *TracyGo) RequestIDKey() string {
	return t.requestID
}

// CorrelationIDFromContext returns the correlation id from the given context, or an empty string.
func (t *TracyGo) CorrelationIDFromContext(ctx context.Context) string {
	if ctx != nil {
		if correlationID, ok := ctx.Value(t.correlationID).(string); ok {
			return correlationID
		}
	}

	return ""
}

// EnsureCorrelationID makes sure the given context has a correlation id set.
func (t *TracyGo) EnsureCorrelationID(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Value(t.correlationID).(string); ok {
		return ctx
	}

	return t.NewContextWithCorrelationID(ctx, uuid.NewString())
}

// RequestIDFromContext returns the requestID from the given context, or an empty string.
func (t *TracyGo) RequestIDFromContext(ctx context.Context) string {
	if ctx != nil {
		if requestID, ok := ctx.Value(t.requestID).(string); ok {
			return requestID
		}
	}

	return ""
}

// NewContextWithCorrelationID sets the correlation id to use in the given context. If ctx is nil, a new context without value is created.
func (t *TracyGo) NewContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, t.correlationID, correlationID) //nolint:staticcheck,revive // intended use
}

// NewContextWithRequestID sets the request id to use in the given context. If ctx is nil, a new context without value is created.
func (t *TracyGo) NewContextWithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, t.requestID, requestID) //nolint:staticcheck,revive // intended use
}
