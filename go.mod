module github.com/alexeyyakimovich/observability-go

go 1.15

require (
	github.com/apache/thrift v0.14.0 // indirect
	github.com/getsentry/sentry-go v0.9.0
	github.com/google/uuid v1.2.0
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	google.golang.org/api v0.39.0 // indirect
)
