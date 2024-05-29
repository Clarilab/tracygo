package http

import (
	"context"
	"net/http"

	"github.com/Clarilab/tracygo/v2"
	"github.com/google/uuid"
)

// CheckTracingIDs is a middleware for net/http that checks if a correlationID and requestID have been set
// and creates a new one if they have not been set yet.
func CheckTracingIDs(t *tracygo.TracyGo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			correlationID := r.Header.Get(t.CorrelationIDKey())
			requestID := r.Header.Get(t.RequestIDKey())

			if correlationID == "" {
				correlationID = uuid.NewString()
			}

			if requestID == "" {
				requestID = uuid.NewString()
			}

			ctx := context.WithValue(r.Context(), t.RequestIDKey(), requestID) //nolint:staticcheck // intended use
			ctx = context.WithValue(ctx, t.CorrelationIDKey(), correlationID)  //nolint:staticcheck // intended use

			w.Header().Set(t.RequestIDKey(), requestID)
			w.Header().Set(t.CorrelationIDKey(), correlationID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
