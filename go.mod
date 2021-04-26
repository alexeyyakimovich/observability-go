module github.com/alexeyyakimovich/observability-go

go 1.15

require (
	github.com/getsentry/sentry-go v0.9.0
	github.com/google/uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.20.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.20.0
)
