package observability

import (
	"fmt"

	"go.opentelemetry.io/otel/label"
)

func fieldsToLabels(fields map[string]interface{}) []label.KeyValue {
	labels := make([]label.KeyValue, len(fields))

	i := 0

	for key, value := range fields {
		labels[i] = label.Key(key).String(fmt.Sprintf("%v", value))
		i++
	}

	return labels
}
