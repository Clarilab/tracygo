package eventhorizon

import (
	"context"

	"github.com/Clarilab/tracygo/v2"
)

type (
	contextMarshalFunc   func(context.Context, map[string]any)
	contextUnmarshalFunc func(context.Context, map[string]any) context.Context
)

// CheckTracingIDs registers marshal and unmarshal functions for correlationIDs in contextes.
// The parameters marshalFunc and unmarshalFunc are to be filled with the RegisterMarshal/UnmarshalFunc function of your eventhorizon library.
func CheckTracingIDs(tracer *tracygo.TracyGo, registerMarshal func(f contextMarshalFunc), registerUnmarshal func(f contextUnmarshalFunc) context.Context) {
	registerMarshal(func(ctx context.Context, vals map[string]any) {
		if correlationID, ok := ctx.Value(tracer.CorrelationIDKey()).(string); ok {
			vals[tracer.CorrelationIDKey()] = correlationID
		}
	})

	registerUnmarshal(func(ctx context.Context, vals map[string]any) context.Context {
		if correlationID, ok := vals[tracer.CorrelationIDKey()].(string); ok {
			ctx = context.WithValue(ctx, tracer.CorrelationIDKey(), correlationID)
		}

		return ctx
	})
}
