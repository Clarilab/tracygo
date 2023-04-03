package tracygo

import (
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

// AtreugoCheckTracingIDs is a useBefore middleware for atreugo that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func (t *TracyGo) AtreugoCheckTracingIDs(ctx *atreugo.RequestCtx) error {
	correlationID := string(ctx.Request.Header.Peek(t.correlationID))
	requestID := string(ctx.Request.Header.Peek(t.requestID))

	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	// set userValue for resty middleware
	ctx.SetUserValue(t.correlationID, correlationID)

	ctx.Response.Header.Set(t.correlationID, correlationID)
	ctx.Response.Header.Set(t.requestID, requestID)

	return ctx.Next()
}

// RestyCheckTracingIDs is a OnBeforeRequest middleware for resty which check if the context has the tracing ids set.
// If they are set, they should be put into the request headers
func (t *TracyGo) RestyCheckTracingIDs(client *resty.Client, request *resty.Request) error {
	requestCtx, ok := request.Context().(*atreugo.RequestCtx)
	if !ok {
		return nil
	}

	// if the correlationID or requestID are present in the atreugo userValues, put them in the resty request
	correlationID, _ := requestCtx.UserValue(t.correlationID).(string)
	if correlationID != "" {
		request.Header.Set(t.correlationID, correlationID)
	}

	request.Header.Set(t.requestID, uuid.New().String())

	return nil
}

// TracyGo is a struct for the tracy object.
type TracyGo struct {
	correlationID string
	requestID     string
}

// Option is an optional func.
type Option func(tracy *TracyGo)

// CorrelationID returns a function that sets the key for the correlationId header.
func CorrelationID(id string) Option {
	return func(tracy *TracyGo) {
		tracy.correlationID = id
	}
}

// RequestID returns a function that sets the key for the requestId header.
func RequestID(id string) Option {
	return func(tracy *TracyGo) {
		tracy.requestID = id
	}
}

// New creates a new TracyGo object and uses the options on it.
func New(options ...Option) *TracyGo {
	tracy := &TracyGo{
		correlationID: "X-Correlation-ID",
		requestID:     "X-Request-ID",
	}

	for _, option := range options {
		option(tracy)
	}

	return tracy
}
