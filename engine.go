// Package observability for logging, tracing and exporting metrics.
package observability

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type engine struct {
	instanceID      uuid.UUID
	logger          Logger
	tracer          Tracer
	metricsExporter MetricsExporter
}

var instance *engine //nolint:gochecknoglobals // singleton

func getInstance() *engine {
	if instance == nil {
		instance = &engine{instanceID: uuid.New()}
	}

	return instance
}

// SetLogger sets global logger.
func SetLogger(logger Logger, fields map[string]interface{}) {
	engine := getInstance()

	fields["instance ID"] = engine.instanceID.String()

	if logger != nil {
		engine.logger = logger.WithFields(fields)
	} else {
		engine.logger = DefaultLogger().WithFields(fields)
	}
}

// GetLogger returns global Logger.
func GetLogger() Logger {
	logger := getInstance().logger

	if logger == nil {
		logger = DefaultLogger()
		SetLogger(logger, map[string]interface{}{})
	}

	return logger
}

// SetTracer sets global tracer.
func SetTracer(tracer Tracer) {
	getInstance().tracer = tracer
}

func getTracer() Tracer {
	return getInstance().tracer
}

// SetMetricsExporter sets global metrics exporter.
func SetMetricsExporter(exporter MetricsExporter) {
	getInstance().metricsExporter = exporter
}

// GetMetricsExporter returns global MetricsExporter or nil, if it wasn't initialized.
func GetMetricsExporter() MetricsExporter {
	return getInstance().metricsExporter
}

// StartOperation logs info message, starts tracing span and gets initial info for gathering metrics.
func StartOperation(ctx context.Context, operationID string, fields map[string]interface{}) (context.Context, Operation) {
	id := uuid.New()

	fields["ID"] = id.String()
	fields["operation"] = operationID

	if logger := GetLogger(); logger != nil {
		logger.WithFields(fields).Info("operation started")
	}

	spanCtx := ctx

	var span Span = nil

	if tracer := getTracer(); tracer != nil {
		spanCtx, span = tracer.Start(ctx, operationID, fields)
	}

	return spanCtx, &operation{
		id:          id,
		span:        span,
		operationID: operationID,
		startedAt:   time.Now(),
		fields:      fields,
	}
}
