package observability

import (
	"net/http"
	"time"
)

// MetricsConfiguration sets metrics collection interval.
type MetricsConfiguration struct {
	CollectionInterval int
}

// MetricsExporter interface.
type MetricsExporter interface {
	HTTPHandler(w http.ResponseWriter, r *http.Request)
	AddOperationCall(operationID string, duration time.Duration, result interface{}, fields map[string]interface{})
	Start()
	Stop()
}
