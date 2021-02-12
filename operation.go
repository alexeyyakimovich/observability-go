package observability

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Operation contains all info to log, trace and gather metrics about concrete operation.
type Operation interface {
	End(result interface{})
}

type operation struct {
	id          uuid.UUID
	startedAt   time.Time
	operationID string
	span        Span
	fields      map[string]interface{}
}

// End logs operation result, closes the span & sends it to tracer and adds call metrics to exporter.
func (op *operation) End(result interface{}) {
	go func() {
		duration := time.Since(op.startedAt)

		// log
		op.fields["result"] = fmt.Sprintf("%v", result)

		logger := GetLogger()
		if logger != nil {
			logger.WithFields(op.fields).Info("operation ended")
		}

		// trace
		if op.span != nil {
			op.span.End()
		}

		// metric
		exporter := getMetricsExporter()
		if exporter != nil {
			exporter.AddOperationCall(op.operationID, duration, result, op.fields)
		}
	}()
}
