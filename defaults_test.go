package observability_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexeyyakimovich/observability-go"
)

const sentryDSN = "https://eccf69905ad948ad8d40716bdd69267f@o477839.ingest.sentry.io/5635267"

var emptyMap = map[string]interface{}{} //nolint:gochecknoglobals // it's a test

func TestDefaults(t *testing.T) {
	t.Parallel()

	observability.InitDefaults("test", "v1.0.0", "tracer:6935", sentryDSN, emptyMap)

	logger := observability.GetLogger()
	require.NotNil(t, logger)

	_, op := observability.StartOperation(context.Background(), "test", emptyMap)

	require.NotNil(t, op)

	op.End("Test OK")
	logger.WithField("time", time.Now().UTC()).Error(errors.New("Test error")) //nolint:goerr113 // it's a test
}
