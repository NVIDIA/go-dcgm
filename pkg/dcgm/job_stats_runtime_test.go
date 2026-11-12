package dcgm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJobStatsAfterFailedInit(t *testing.T) {
	requireLoadedLibraryWithoutSession(t)
	cases := []struct {
		name, symbol string
		call         func() error
	}{
		{"watch", "dcgmWatchJobFields", func() error { return WatchJobFields(GroupHandle{}, time.Second, time.Minute, 0) }},
		{"start", "dcgmJobStartStats", func() error { return JobStart(GroupHandle{}, "valid-job") }},
		{"stop", "dcgmJobStopStats", func() error { return JobStop("valid-job") }},
		{"get", "dcgmJobGetStats", func() error { _, err := JobGetStats("valid-job"); return err }},
		{"remove", "dcgmJobRemove", func() error { return JobRemove("valid-job") }},
		{"remove all", "dcgmJobRemoveAll", JobRemoveAll},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.True(t, dcgmSymbolAvailable(tc.symbol), "regression requires a loaded native symbol")
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, tc.call(), &unavailable)
			require.Equal(t, tc.symbol, unavailable.Symbol)
		})
	}
}
