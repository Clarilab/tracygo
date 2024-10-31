package eventhorizon

import (
	"context"

	"github.com/Clarilab/tracygo/v2"
)

// CheckTracingIDs registers marshal and unmarshal functions for correlationIDs in contextes.
// The parameter functions registerMarshal and registerUnmarshal are supposed to be filled with functions,
// that wrap around your respective eventHorizon libraries call to RegisterContextMarshaler/RegisterContextUnmarshaler
func CheckTracingIDs(
	tracer *tracygo.TracyGo,
	registerMarshal func(func(context.Context, map[string]any)),
	registerUnmarshal func(func(context.Context, map[string]any) context.Context),
) {
	marshalFunc := func(ctx context.Context, vals map[string]any) {
		if correlationID, ok := ctx.Value(tracer.CorrelationIDKey()).(string); ok {
			vals[tracer.CorrelationIDKey()] = correlationID
		}
	}

	registerMarshal(marshalFunc)

	unmarshalFunc := func(ctx context.Context, vals map[string]any) context.Context {
		if correlationID, ok := vals[tracer.CorrelationIDKey()].(string); ok {
			ctx = context.WithValue(ctx, tracer.CorrelationIDKey(), correlationID)
		}

		return ctx
	}

	registerUnmarshal(unmarshalFunc)
}
