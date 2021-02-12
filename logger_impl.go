package observability

import log "github.com/sirupsen/logrus"

type logrusWrapper struct {
	logrus *log.Logger
	fields map[string]interface{}
}

func (wrapper logrusWrapper) Debug(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Debug(args...)
}

func (wrapper logrusWrapper) Debugf(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Debugf(format, args...)
}

func (wrapper logrusWrapper) Info(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Info(args...)
}

func (wrapper logrusWrapper) Infof(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Infof(format, args...)
}

func (wrapper logrusWrapper) Error(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Error(args...)
}

func (wrapper logrusWrapper) Errorf(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Errorf(format, args...)
}

func (wrapper logrusWrapper) Warning(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Warning(args...)
}

func (wrapper logrusWrapper) Warningf(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Warningf(format, args...)
}

func (wrapper logrusWrapper) Panic(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Panic(args...)
}

func (wrapper logrusWrapper) Panicf(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Panicf(format, args...)
}

func (wrapper logrusWrapper) Fatal(args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Fatal(args...)
}

func (wrapper logrusWrapper) Fatalf(format string, args ...interface{}) {
	wrapper.logrus.WithFields(wrapper.fields).Fatalf(format, args...)
}

func (wrapper logrusWrapper) WithFields(fields map[string]interface{}) Logger {
	f := wrapper.fields
	for k, v := range fields {
		f[k] = v
	}

	return &logrusWrapper{
		logrus: wrapper.logrus,
		fields: f,
	}
}

func (wrapper logrusWrapper) WithField(key string, value interface{}) Logger {
	f := wrapper.fields
	f[key] = value

	return &logrusWrapper{
		logrus: wrapper.logrus,
		fields: f,
	}
}

func (wrapper logrusWrapper) SetLevel(level Level) {
	wrapper.logrus.SetLevel(log.Level(level))
}

// NewLogger creates new Logger instance.
func NewLogger(config LoggerConfiguration, appID, version string, fields map[string]interface{}) Logger {
	wrapper := logrusWrapper{
		logrus: log.New(),
		fields: fields,
	}

	wrapper.SetLevel(config.MinLevel)

	if config.SentryDSN != "" {
		sentryHook, err := NewSentryHook(config.SentryDSN, appID+"@"+version)
		if err != nil {
			wrapper.Error(err)
		} else {
			wrapper.logrus.Hooks.Add(sentryHook)
		}
	}

	return wrapper
}
