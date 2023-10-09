package tracygo_test

import (
	"context"
	"testing"

	"github.com/Clarilab/tracygo/v2"
)

func Test_NewContext(t *testing.T) {
	tracer := tracygo.New()

	t.Run("with context", func(t *testing.T) {

		ctx := tracer.NewContext(context.Background(), "Zitronenbaum")

		id, ok := ctx.Value(tracygo.CorrelationID("X-Correlation-ID")).(string)
		if !ok {
			t.Error("invalid type for correlationID")
		}

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("nil context", func(t *testing.T) {

		ctx := tracer.NewContext(nil, "Zitronenbaum") //nolint:staticcheck // intended use for testing

		id, ok := ctx.Value(tracygo.CorrelationID("X-Correlation-ID")).(string)
		if !ok {
			t.Error("invalid type for correlationID")
		}

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})
}

func Test_FromContext(t *testing.T) {
	tracer := tracygo.New()

	t.Run("correlationID exists", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), tracygo.CorrelationID("X-Correlation-ID"), "Zitronenbaum")

		id := tracer.FromContext(ctx)

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("correlationID does not exist", func(t *testing.T) {
		id := tracer.FromContext(context.Background())

		if id != "" {
			t.Errorf("expected '', got '%s'", id)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		id := tracer.FromContext(nil) //nolint:staticcheck // intended use for testing

		if id != "" {
			t.Errorf("expected '', got '%s'", id)
		}
	})
}

func Test_GetIDs(t *testing.T) {
	tracer := tracygo.New()

	t.Run("correlationID type", func(t *testing.T) {
		key := tracer.GetCorrelationID()

		if key != "X-Correlation-ID" {
			t.Errorf("expected 'X-Correlation-ID', got '%s'", key)
		}
	})

	t.Run("correlationID string", func(t *testing.T) {
		key := tracer.GetCorrelationID()

		if key.String() != "X-Correlation-ID" {
			t.Errorf("expected 'X-Correlation-ID', got '%s'", key)
		}
	})

	t.Run("requestID type", func(t *testing.T) {
		key := tracer.GetRequestID()

		if key != "X-Request-ID" {
			t.Errorf("expected 'X-Request-ID', got '%s'", key)
		}
	})

	t.Run("requestID string", func(t *testing.T) {
		key := tracer.GetRequestID()

		if key.String() != "X-Request-ID" {
			t.Errorf("expected 'X-Request-ID', got '%s'", key)
		}
	})
}

func Test_Options(t *testing.T) {
	tracer := tracygo.New(
		tracygo.WithCorrelationID("X-Correlation-Zitronenbaum"),
		tracygo.WithRequestID("X-Request-Zitronenbaum"),
	)

	t.Run("correlationID", func(t *testing.T) {
		key := tracer.GetCorrelationID()

		if key.String() != "X-Correlation-Zitronenbaum" {
			t.Errorf("expected 'X-Correlation-Zitronenbaum', got '%s'", key)
		}
	})

	t.Run("requestID", func(t *testing.T) {
		key := tracer.GetRequestID()

		if key.String() != "X-Request-Zitronenbaum" {
			t.Errorf("expected 'X-Request-Zitronenbaum', got '%s'", key)
		}
	})
}
