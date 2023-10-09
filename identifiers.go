package tracygo

// CorrelationID is the key for the correlationID.
type CorrelationID string

// String returns the string representation of the correlationID.
func (c CorrelationID) String() string {
	return string(c)
}

// RequestID is the key for the requestID.
type RequestID string

// String returns the string representation of the RequestID.
func (c RequestID) String() string {
	return string(c)
}
