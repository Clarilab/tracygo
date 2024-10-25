package grpc

import (
	"context"
	"strings"

	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CheckTracingIDs is a middleware for grpc that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		var correlationID string
		var requestID string

		if values := md[strings.ToLower(string(t.CorrelationIDKey()))]; len(values) == 1 {
			correlationID = values[0]
		}

		if values := md[strings.ToLower(string(t.RequestIDKey()))]; len(values) == 1 {
			requestID = values[0]
		}

		if correlationID == "" {
			correlationID = uuid.NewString()

			md.Append(string(t.CorrelationIDKey()), correlationID)
		}

		if requestID == "" {
			requestID = uuid.NewString()

			md.Append(string(t.RequestIDKey()), requestID)
		}

		if err := grpc.SetTrailer(ctx, md); err != nil {
			return handler(ctx, req)
		}

		ctx = metadata.AppendToOutgoingContext(
			ctx,
			string(t.CorrelationIDKey()), correlationID,
			string(t.RequestIDKey()), requestID,
		)

		ctx = context.WithValue(ctx, t.CorrelationIDKey(), correlationID)
		ctx = context.WithValue(ctx, t.RequestIDKey(), requestID)

		return handler(ctx, req)
	}
}
