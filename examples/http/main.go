package main

import (
	"net/http"

	httptracygo "github.com/Clarilab/tracygo/middleware/http/v2"
	restytracygo "github.com/Clarilab/tracygo/middleware/resty/v2"

	"github.com/Clarilab/tracygo/v2"
	"github.com/go-resty/resty/v2"
)

func main() {
	tracer := tracygo.New()

	// setup separate server for main server to call
	setupSeparateServer()
	setupMainServer(tracer)

	// request to main server
	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add(tracer.CorrelationIDKey(), "Zitronenbaum")

	// call main server
	resp, err := new(http.Client).Do(request)
	if err != nil {
		panic(err)
	}

	if resp.Header.Get(tracer.CorrelationIDKey()) != "Zitronenbaum" {
		panic("X-Correlation-ID header is not set in response")
	}

	if resp.Header.Get(tracer.RequestIDKey()) == "" {
		panic("X-Request-ID header is not set in response")
	}
}

func setupMainServer(tracer *tracygo.TracyGo) {
	mux := http.NewServeMux()

	restyClient := resty.New()
	restyClient.OnBeforeRequest(restytracygo.RestyCheckTracingIDs(tracer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(tracer.CorrelationIDKey()) != "Zitronenbaum" {
			panic("X-Correlation-ID header is not set in context")
		}

		correlationID, ok := r.Context().Value(tracer.CorrelationIDKey()).(string)
		if !ok || correlationID != "Zitronenbaum" {
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

	go func() { _ = http.ListenAndServe(":8080", server) }()
}

func setupSeparateServer() {
	const (
		headerCorrelationID = "X-Correlation-ID"
		headerRequestID     = "X-Request-ID"
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerCorrelationID) != "Zitronenbaum" {
			panic("X-Correlation-ID header is not set")
		}

		if r.Header.Get(headerRequestID) == "" {
			panic("X-Request-ID header is not set")
		}

		w.WriteHeader(http.StatusOK)
	})

	go func() { _ = http.ListenAndServe(":8081", mux) }()
}
