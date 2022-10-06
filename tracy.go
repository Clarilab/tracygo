package tracygo

import (
	"errors"

	"context"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"

	clarilab "github.com/Clarilab/eventhorizon"
)

// CheckRequestID is a useBefore middleware that checks if a correlationID is set if that is not the case it creates a new one
// and creates a new requestID which both get written into the userValues
func (t *TracyGo) CheckRequestID(ctx *atreugo.RequestCtx) error {
	correlationId := string(ctx.Request.Header.Peek(t.correlationId))

	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	requestId := uuid.New().String()

	ctx.SetUserValue(t.correlationId, correlationId)
	ctx.SetUserValue(t.requestId, requestId)

	return ctx.Next()
}

// WriteHeader is a useAfter middle which takes the correlationId and requestId and writes then into the response Header.
func (t *TracyGo) WriteHeader(ctx *atreugo.RequestCtx) error {
	correlationId := ctx.UserValue(t.correlationId).(string)
	requestId := ctx.UserValue(t.requestId).(string)

	ctx.Response.Header.Set(t.correlationId, correlationId)
	ctx.Response.Header.Set(t.requestId, requestId)

	return ctx.Next()
}

// CheckTracingIDs is a OnBeforeRequest middleware which check if the context has the tracing ids set.
// If they are set, they should be put into the request headers
func (t *TracyGo) CheckTracingIDs(client *resty.Client, request *resty.Request) error {
	correlationId, ok := request.Context().Value(t.correlationId).(string)
	if !ok {
		return errors.New("correlationId not found")
	}

	requestId, ok := request.Context().Value(t.requestId).(string)
	if !ok {
		return errors.New("requestId not found")
	}

	request.Header.Set(t.correlationId, correlationId)
	request.Header.Set(t.requestId, requestId)

	return nil
}

// TracyGo is a struct for the tracy object.
type TracyGo struct {
	correlationId string
	requestId     string
}

// Option is an optional func.
type Option func(tracy *TracyGo)

// CorrelationId returns a function that sets the correlationId of the header.
func CorrelationId(id string) Option {
	return func(tracy *TracyGo) {
		tracy.correlationId = id
	}
}

// RequestId returns a function that sets the requestId of the header.
func RequestId(id string) Option {
	return func(tracy *TracyGo) {
		tracy.requestId = id
	}
}

// New creates a new Tracygo object and uses the options on it.
func New(options ...Option) *TracyGo {
	tracy := &TracyGo{
		correlationId: "X-Correlation-ID",
		requestId:     "X-Request-ID",
	}

	for _, option := range options {
		option(tracy)
	}

	return tracy
}

// DefaultNamespace is the namespace to use if not set in the context.
const DefaultNamespace = "default"

// Strings used to marshal context values.
const (
	namespaceKeyStr = "clairlab_namespace"
)

func init() {
	clarilab.RegisterContextMarshaler(func(ctx context.Context, vals map[string]interface{}) {
		if ns, ok := ctx.Value(namespaceKey).(string); ok {
			vals[namespaceKeyStr] = ns
		}
	})
	clarilab.RegisterContextUnmarshaler(func(ctx context.Context, vals map[string]interface{}) context.Context {
		if ns, ok := vals[namespaceKeyStr].(string); ok {
			ctx = NewContext(ctx, ns)
		}

		return ctx
	})
}

type contextKey int

const (
	namespaceKey contextKey = iota
)

// FromContext returns the namespace from the context, or the default namespace.
func FromContext(ctx context.Context) string {
	if ns, ok := ctx.Value(namespaceKey).(string); ok {
		return ns
	}

	return DefaultNamespace
}

// NewContext sets the namespace to use in the context. The namespace is used to
// determine which database to use, among other things.
func NewContext(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceKey, namespace)
}
