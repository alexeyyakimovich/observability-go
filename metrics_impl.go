package observability

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
)

// MetricsExporter contains settings for exporting metrics.
type metricsExporter struct {
	// server
	exporter *prometheus.Exporter

	// operations metrics
	meter   metric.Meter
	counter metric.Int64Counter
	rps     metric.Int64ValueRecorder

	// utils
	mutex   sync.Mutex
	started bool

	// configuration
	config MetricsConfiguration
	fields map[string]interface{}
}

// HTTPHandler to handle Prometheus metrics scrapper requests.
func (exporter *metricsExporter) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if exporter.exporter != nil {
		exporter.exporter.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// AddOperationCall adds duration metric and increases request count by 1 for operation.
func (exporter *metricsExporter) AddOperationCall(
	operationID string,
	duration time.Duration,
	result interface{},
	fields map[string]interface{}) {
	labels := fieldsToLabels(fields)

	labels = append(labels,
		label.Key("result").String(fmt.Sprintf("%v", result)),
		label.Key("operation").String(operationID))

	exporter.counter.Add(context.Background(), 1, labels...)
	exporter.rps.Record(context.Background(), duration.Nanoseconds(), labels...)
}

func (exporter *metricsExporter) updateServiceMetrics() {
	mem, _ := exporter.meter.NewInt64ValueRecorder("mem_usage",
		metric.WithDescription("Amount of memory used."),
	)
	used, _ := exporter.meter.NewFloat64ValueRecorder("disk_usage",
		metric.WithDescription("Amount of disk used."),
	)
	quota, _ := exporter.meter.NewFloat64ValueRecorder("disk_quota",
		metric.WithDescription("Amount of disk quota available."),
	)
	goroutines, _ := exporter.meter.NewInt64ValueRecorder("num_goroutines",
		metric.WithDescription("Amount of goroutines running."),
	)

	var m runtime.MemStats

	for {
		runtime.ReadMemStats(&m)

		var stat syscall.Statfs_t

		wd, _ := os.Getwd()

		err := syscall.Statfs(wd, &stat)
		if err != nil {
			GetLogger().WithField("scope", "metric exporter").Error(err)
		}

		all := float64(stat.Blocks) * float64(stat.Bsize)
		free := float64(stat.Bfree) * float64(stat.Bsize)

		exporter.meter.RecordBatch(context.Background(),
			fieldsToLabels(exporter.fields),
			used.Measurement(all-free),
			quota.Measurement(all),
			mem.Measurement(int64(m.Sys)),
			goroutines.Measurement(int64(runtime.NumGoroutine())),
		)
		time.Sleep(time.Second * time.Duration(exporter.config.CollectionInterval))

		if !exporter.started {
			break
		}
	}
}

// Reset stops metrics collection routine.
func (exporter *metricsExporter) Stop() {
	exporter.started = false
}

func (exporter *metricsExporter) Start() {
	exporter.mutex.Lock()
	defer exporter.mutex.Unlock()

	if !exporter.started {
		exporter.started = true

		go exporter.updateServiceMetrics()
	}
}

// NewMetricsExporter initializes new prometheus exporter.
func NewMetricsExporter(config MetricsConfiguration, fields map[string]interface{}) MetricsExporter {
	if config.CollectionInterval <= 0 {
		return nil
	}

	const scope = "metric exporter"

	log := GetLogger()

	// NOTE: will be moved back to go.opentelemetry.io/otel/metric/global
	meter := otel.GetMeterProvider().Meter("github.com/alexeyyakimovich/observability-go")

	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.WithField("scope", scope).Error(err)
	}

	counter, err := meter.NewInt64Counter("operation_requests_total",
		metric.WithDescription("Number of requests"))
	if err != nil {
		log.WithField("scope", scope).Error(err)
	}

	rps, err := meter.NewInt64ValueRecorder("requests_duration",
		metric.WithDescription("Request duration."),
	)
	if err != nil {
		log.WithField("scope", scope).Error(err)
	}

	result := metricsExporter{
		exporter: exporter,
		fields:   fields,
		meter:    meter,
		counter:  counter,
		rps:      rps,
		started:  false,
		mutex:    sync.Mutex{},
		config:   config,
	}

	return &result
}
