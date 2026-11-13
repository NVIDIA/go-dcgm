//go:build linux && cgo

package dcgm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigurationAfterFailedInit(t *testing.T) {
	requireLoadedLibraryWithoutSession(t)
	require.True(t, dcgmSymbolAvailable("dcgmConfigGet"))
	calls := map[string]func() error{
		"get":     func() error { _, err := ConfigGet(GroupHandle{}, ConfigCurrentState); return err },
		"set":     func() error { return ConfigSet(GroupHandle{}, GPUConfig{}) },
		"enforce": func() error { return ConfigEnforce(GroupHandle{}) },
		"reset":   func() error { return ConfigReset(GroupHandle{}) },
		"profile set": func() error {
			return ConfigSetWorkloadPowerProfile(GroupHandle{}, WorkloadProfileSet, PowerProfileMask{1})
		},
		"profile info": func() error { _, err := GetDeviceWorkloadPowerProfileInfo(0); return err },
		"profile name": func() error { _, err := PowerProfileName(PowerProfileCompute); return err },
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, call(), &unavailable)
		})
	}
}
