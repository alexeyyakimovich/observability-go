package observability

import "context"

// TracingConfiguration sets tracer service endpoint and sampling rate.
type TracingConfiguration struct {
	TracerEndpoint      string
	SamplingProbability float64
}

// Span interface for wrapping Jaeger span.
type Span interface {
	End()
}

// Tracer interface for wrapping Jaeger tracer.
type Tracer interface {
	// Starts a span.
	Start(ctx context.Context, operationID string, fields map[string]interface{}) (context.Context, Span)
}
