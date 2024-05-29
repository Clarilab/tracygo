package tracygo

// Option is an optional func.
type Option func(tracy *TracyGo)

// WithCorrelationID returns a function that sets the key for the correlation id header.
func WithCorrelationID(id string) Option {
	return func(tracy *TracyGo) {
		tracy.correlationID = id
	}
}

// WithRequestID returns a function that sets the key for the request id header.
func WithRequestID(id string) Option {
	return func(tracy *TracyGo) {
		tracy.requestID = id
	}
}
