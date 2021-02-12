package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// NewTracer traces exporter with provided values.
func NewTracer(config TracingConfiguration, appID string, log Logger) Tracer {
	wrapper := &tracerWrapper{}

	if config.SamplingProbability > 0 && config.TracerEndpoint != "" {
		// NOTE: will be moved back to go.opentelemetry.io/otel/api/global.
		wrapper.tracer = otel.GetTracerProvider().Tracer("github.com/alexeyyakimovich/observability-go")

		configureJaegerTracer(config, appID, log)
	}

	return wrapper
}

func configureJaegerTracer(config TracingConfiguration, appID string, log Logger) {
	sdkConfig := sdktrace.Config{ //nolint:exhaustivestruct // the only options we need.
		DefaultSampler: sdktrace.TraceIDRatioBased(config.SamplingProbability),
	}

	jExporter, err := jaeger.NewRawExporter(
		jaeger.WithAgentEndpoint(config.TracerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: appID,
			Tags:        []label.KeyValue{},
		}),
	)
	if err != nil {
		log.WithField("scope", "tracer").Fatal(err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdkConfig),
		sdktrace.WithSyncer(jExporter),
	)

	// NOTE: will be moved back to go.opentelemetry.io/otel/api/global.
	otel.SetTracerProvider(tracerProvider)
}

type spanWrapper struct {
	span trace.Span
}

func (wrapper spanWrapper) End() {
	if wrapper.span != nil {
		wrapper.span.End()
	}
}

type tracerWrapper struct {
	tracer trace.Tracer
}

func (wrapper tracerWrapper) Start(ctx context.Context, operationID string, fields map[string]interface{}) (context.Context, Span) {
	if wrapper.tracer != nil {
		fields["operation"] = operationID

		labels := fieldsToLabels(fields)

		opts := []trace.SpanOption{
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(labels...),
		}
		spanCtx, span := wrapper.tracer.Start(ctx, operationID, opts...)

		wr := spanWrapper{span: span}

		return spanCtx, wr
	}

	GetLogger().WithFields(map[string]interface{}{"scope": "tracer"}).Warning("tracer wasn't configured")

	return ctx, spanWrapper{}
}
