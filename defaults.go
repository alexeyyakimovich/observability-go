package observability

const (
	// DefaultMetricsCollectionInterval is 15 seconds.
	DefaultMetricsCollectionInterval = 15
	// DefaultLogLevel is ErrorLevel.
	DefaultLogLevel = ErrorLevel
	// DefaultSamplingRate is 5%.
	DefaultSamplingRate = 0.05
)

// InitDefaults initializes Logger, MetricsExporter and Tracer with default values.
func InitDefaults(appID, version, tracerEndpoint, sentryDSN string, fields map[string]interface{}) {
	fields["app id"] = appID
	fields["version"] = version

	// Logger setup
	loggerConfig := LoggerConfiguration{
		MinLevel:  DefaultLogLevel,
		SentryDSN: sentryDSN,
	}
	logger := NewLogger(loggerConfig, appID, version, fields)

	SetLogger(logger, fields)

	// Exporter setup
	exporterConfig := MetricsConfiguration{
		CollectionInterval: DefaultMetricsCollectionInterval,
	}
	exporter := NewMetricsExporter(exporterConfig, logger, fields)

	exporter.Start()
	SetMetricsExporter(exporter)

	// Tracer sertup
	if tracerEndpoint != "" {
		tracerConfig := TracingConfiguration{
			SamplingProbability: DefaultSamplingRate,
			TracerEndpoint:      tracerEndpoint,
		}
		tracer := NewTracer(tracerConfig, appID, logger)

		SetTracer(tracer)
	}
}
