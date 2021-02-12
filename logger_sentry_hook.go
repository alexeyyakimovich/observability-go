package observability

import (
	"fmt"
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

const flushTimeout = 2

// severityMap is a mapping of logrus log level to sentry log level.
var severityMap = map[log.Level]sentry.Level{ //nolint:gochecknoglobals // application-level mapping
	log.DebugLevel: sentry.LevelDebug,
	log.InfoLevel:  sentry.LevelInfo,
	log.WarnLevel:  sentry.LevelWarning,
	log.ErrorLevel: sentry.LevelError,
	log.FatalLevel: sentry.LevelFatal,
	log.PanicLevel: sentry.LevelFatal,
}

// SentryHook implements logrus.Hook to send errors to sentry.
type SentryHook struct {
	client *sentry.Client
	levels []log.Level
}

// SentryEventIdentityModifier is a sentry event modifier that simply passes
// through the event.
type SentryEventIdentityModifier struct{}

// ApplyToEvent simply returns the event (ignoring the hint).
func (m *SentryEventIdentityModifier) ApplyToEvent(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	return event
}

// NewSentryHook creates a sentry hook for logrus given a sentry dsn.
func NewSentryHook(dsn, release string) (*SentryHook, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{ //nolint:exhaustivestruct // this is exact fieldset we need.
		Dsn:     dsn,
		Release: release,
	})
	if err != nil {
		log.WithField("error", err).Error("Unable to initialize Sentry")

		return nil, fmt.Errorf("sentry initialization error: %w", err)
	}

	return &SentryHook{
		client: client,
		levels: []log.Level{
			log.WarnLevel,
			log.ErrorLevel,
			log.FatalLevel,
			log.PanicLevel,
		},
	}, nil
}

// Levels returns the levels this hook is enabled for. This is a part
// of logrus.Hook.
func (h *SentryHook) Levels() []log.Level {
	return h.levels
}

// Fire is an event handler for logrus. This is a part of logrus.Hook.
func (h *SentryHook) Fire(e *log.Entry) error {
	event := sentry.NewEvent()
	event.Message = e.Message
	event.Level = severityMap[e.Level]
	event.Timestamp = e.Time.UTC()

	var err error

	for k, v := range e.Data {
		if k == log.ErrorKey {
			err = v.(error)
		} else {
			event.Extra[k] = v
		}
	}

	if err != nil {
		stacktrace := sentry.ExtractStacktrace(err)
		event.Exception = []sentry.Exception{{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: stacktrace,
		}}
	}

	sentryModifier := &SentryEventIdentityModifier{}

	h.client.CaptureEvent(event, nil, sentryModifier)

	h.client.Flush(flushTimeout * time.Second)

	return nil
}
