package tracygo_test

import (
	"context"
	"testing"

	"github.com/Clarilab/tracygo/v2"
)

const (
	keyCorrelationID = "X-Correlation-ID"
	keyRequestID     = "X-Request-ID"
)

func Test_NewContext(t *testing.T) {
	t.Parallel()

	tracer := tracygo.New()

	t.Run("with correlationID", func(t *testing.T) {
		t.Parallel()

		ctx := tracer.NewContextWithCorrelationID(context.Background(), "Zitronenbaum")

		id, ok := ctx.Value(keyCorrelationID).(string)
		if !ok {
			t.Error("invalid type for correlationID")
		}

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("with requestID", func(t *testing.T) {
		t.Parallel()

		ctx := tracer.NewContextWithRequestID(context.Background(), "Zitronenbaum")

		id, ok := ctx.Value(keyRequestID).(string)
		if !ok {
			t.Error("invalid type for requestID")
		}

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		t.Parallel()

		ctx := tracer.NewContextWithCorrelationID(nil, "Zitronenbaum") //nolint:staticcheck // intended use for testing

		id, ok := ctx.Value("X-Correlation-ID").(string)
		if !ok {
			t.Error("invalid type for correlationID")
		}

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})
}

func Test_FromContext(t *testing.T) {
	t.Parallel()

	tracer := tracygo.New()

	t.Run("correlationID exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), keyCorrelationID, "Zitronenbaum") //nolint:staticcheck // intended use for testing

		id := tracer.CorrelationIDromContext(ctx)

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("correlationID does not exist", func(t *testing.T) {
		t.Parallel()

		id := tracer.CorrelationIDromContext(context.Background())

		if id != "" {
			t.Errorf("expected '', got '%s'", id)
		}
	})

	t.Run("requestID exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), keyRequestID, "Zitronenbaum") //nolint:staticcheck // intended use for testing

		id := tracer.RequestIDFromContext(ctx)

		if id != "Zitronenbaum" {
			t.Errorf("expected 'Zitronenbaum', got '%s'", id)
		}
	})

	t.Run("requestID does not exist", func(t *testing.T) {
		t.Parallel()

		id := tracer.RequestIDFromContext(context.Background())

		if id != "" {
			t.Errorf("expected '', got '%s'", id)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		t.Parallel()

		id := tracer.CorrelationIDromContext(nil) //nolint:staticcheck // intended use for testing

		if id != "" {
			t.Errorf("expected '', got '%s'", id)
		}
	})
}

func Test_GetIDs(t *testing.T) {
	t.Parallel()

	tracer := tracygo.New()

	t.Run("correlationID type", func(t *testing.T) {
		t.Parallel()

		key := tracer.CorrelationIDKey()

		if key != "X-Correlation-ID" {
			t.Errorf("expected 'X-Correlation-ID', got '%s'", key)
		}
	})

	t.Run("correlationID string", func(t *testing.T) {
		t.Parallel()

		key := tracer.CorrelationIDKey()

		if key != "X-Correlation-ID" {
			t.Errorf("expected 'X-Correlation-ID', got '%s'", key)
		}
	})

	t.Run("requestID type", func(t *testing.T) {
		t.Parallel()

		key := tracer.RequestIDKey()

		if key != "X-Request-ID" {
			t.Errorf("expected 'X-Request-ID', got '%s'", key)
		}
	})

	t.Run("requestID string", func(t *testing.T) {
		t.Parallel()

		key := tracer.RequestIDKey()

		if key != "X-Request-ID" {
			t.Errorf("expected 'X-Request-ID', got '%s'", key)
		}
	})
}

func Test_Options(t *testing.T) {
	t.Parallel()

	tracer := tracygo.New(
		tracygo.WithCorrelationID("X-Correlation-Zitronenbaum"),
		tracygo.WithRequestID("X-Request-Zitronenbaum"),
	)

	t.Run("correlationID", func(t *testing.T) {
		t.Parallel()

		key := tracer.CorrelationIDKey()

		if key != "X-Correlation-Zitronenbaum" {
			t.Errorf("expected 'X-Correlation-Zitronenbaum', got '%s'", key)
		}
	})

	t.Run("requestID", func(t *testing.T) {
		t.Parallel()

		key := tracer.RequestIDKey()

		if key != "X-Request-Zitronenbaum" {
			t.Errorf("expected 'X-Request-Zitronenbaum', got '%s'", key)
		}
	})
}
