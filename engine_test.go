package observability_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexeyyakimovich/observability-go"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := observability.GetLogger()

	require.NotNil(t, logger)

	observability.SetLogger(nil, emptyMap)

	logger = observability.GetLogger()

	require.NotNil(t, logger)
}
