package resty

import (
	"github.com/Clarilab/tracygo/v2"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// RestyCheckTracingIDs is a OnBeforeRequest middleware for resty which check if the context has the tracing ids set.
// If they are set, they should be put into the request headers.
func RestyCheckTracingIDs(t *tracygo.TracyGo) func(client *resty.Client, request *resty.Request) error {
	return func(client *resty.Client, request *resty.Request) error {
		request.Header.Set(t.RequestIDKey(), uuid.NewString())

		correlationID, ok := request.Context().Value(t.CorrelationIDKey()).(string)
		if ok && correlationID != "" {
			request.Header.Set(t.CorrelationIDKey(), correlationID)

			return nil
		}

		request.Header.Set(t.CorrelationIDKey(), uuid.NewString())

		return nil
	}
}
