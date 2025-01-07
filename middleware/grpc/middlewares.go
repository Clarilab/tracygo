package grpc

import (
	"context"
	"strings"

	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CheckTracingIDs is a middleware for grpc servers that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var correlationID string
		var requestID string

		if values := md[strings.ToLower(t.CorrelationIDKey())]; len(values) == 1 {
			correlationID = values[0]
		}

		if values := md[strings.ToLower(t.RequestIDKey())]; len(values) == 1 {
			requestID = values[0]
		}

		if correlationID == "" {
			correlationID = uuid.NewString()

			md.Append(t.CorrelationIDKey(), correlationID)
		}

		if requestID == "" {
			requestID = uuid.NewString()

			md.Append(t.RequestIDKey(), requestID)
		}

		if err := grpc.SetTrailer(ctx, md); err != nil {
			return handler(ctx, req)
		}

		ctx = metadata.AppendToOutgoingContext(
			ctx,
			t.CorrelationIDKey(), correlationID,
			t.RequestIDKey(), requestID,
		)

		ctx = context.WithValue(ctx, t.CorrelationIDKey(), correlationID) //nolint:staticcheck // intended use
		ctx = context.WithValue(ctx, t.RequestIDKey(), requestID)         //nolint:staticcheck // intended use

		return handler(ctx, req)
	}
}

// ClientMiddleware is a middleware for grpc clients. It ensures that a correlationID is set.
func ClientMiddleware(t *tracygo.TracyGo) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = t.EnsureCorrelationID(ctx)
		cID := t.CorrelationIDFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, t.CorrelationIDKey(), cID)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
