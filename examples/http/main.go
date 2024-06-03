package main

import (
	"context"
	"net/http"

	"github.com/Clarilab/tracygo/v2"
	httptracygo "github.com/Clarilab/tracygo/v2/middleware/http"
	restytracygo "github.com/Clarilab/tracygo/v2/middleware/resty"
	"github.com/go-resty/resty/v2"
)

const (
	correlationValue = "Zitronenbaum"
)

func main() {
	tracer := tracygo.New()

	// setup separate server for main server to call
	setupSeparateServer()

	// setup main server
	setupMainServer(tracer)

	resp, err := callMainServer(context.Background(), tracer.CorrelationIDKey())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.Header.Get(tracer.CorrelationIDKey()) != correlationValue {
		panic("X-Correlation-ID header is not set in response")
	}

	if resp.Header.Get(tracer.RequestIDKey()) == "" {
		panic("X-Request-ID header is not set in response")
	}
}

func callMainServer(ctx context.Context, correlationID string) (*http.Response, error) {
	// request to main server
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add(correlationID, correlationValue)

	// call main server
	return new(http.Client).Do(request) //nolint:wrapcheck
}

func setupMainServer(tracer *tracygo.TracyGo) {
	mux := http.NewServeMux()

	restyClient := resty.New()
	restyClient.OnBeforeRequest(restytracygo.CheckTracingIDs(tracer))

	mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		if r.Header.Get(tracer.CorrelationIDKey()) != correlationValue {
			panic("X-Correlation-ID header is not set in context")
		}

		correlationID, ok := r.Context().Value(tracer.CorrelationIDKey()).(string)
		if !ok || correlationID != correlationValue {
			panic("correlation id is not set in context")
		}

		requestID, ok := r.Context().Value(tracer.RequestIDKey()).(string)
		if !ok || requestID == "" {
			panic("request id is not set")
		}

		// call separate server
		_, err := restyClient.R().
			SetContext(r.Context()).
			EnableTrace().
			Get("http://localhost:8081")
		if err != nil {
			panic(err)
		}
	})

	server := httptracygo.CheckTracingIDs(tracer)(mux)

	go func() { panic(http.ListenAndServe(":8080", server)) }()
}

func setupSeparateServer() {
	const (
		headerCorrelationID = "X-Correlation-ID"
		headerRequestID     = "X-Request-ID"
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerCorrelationID) != correlationValue {
			panic("X-Correlation-ID header is not set")
		}

		if r.Header.Get(headerRequestID) == "" {
			panic("X-Request-ID header is not set")
		}

		w.WriteHeader(http.StatusOK)
	})

	go func() { panic(http.ListenAndServe(":8081", mux)) }()
}
