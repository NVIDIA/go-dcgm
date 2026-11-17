package dcgm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// checkAdministrationWithoutSession verifies each public entry point reaches its own session guard.
func checkAdministrationWithoutSession(t *testing.T) {
	t.Helper()
	for _, test := range []struct {
		name, symbol string
		call         func() error
	}{
		{"health", "dcgmHostengineIsHealthy", func() error { _, err := HostengineIsHealthy(); return err }},
		{"logging", "dcgmHostengineSetLoggingSeverity", func() error { return HostengineSetLoggingSeverity(0, LoggingNone) }},
		{"environment", "dcgmHostengineEnvironmentVariableInfo", func() error { _, err := HostengineEnvironmentVariable("DCGM_DIAG_REMOVED"); return err }},
		{"modules", "dcgmModuleGetStatuses", func() error { _, err := ModuleGetStatuses(); return err }},
		{"module name", "dcgmModuleIdToName", func() error { _, err := ModuleIDToName(ModuleHealth); return err }},
		{"denylist", "dcgmModuleDenylist", func() error { return ModuleDenylist(ModuleHealth) }},
		{"groups", "dcgmGroupGetAllIds", func() error { _, err := GroupGetAllIDs(); return err }},
		{"remove GPU", "dcgmGroupRemoveDevice", func() error { return RemoveFromGroup(GroupHandle{}, 0) }},
		{"field group", "dcgmFieldGroupGetInfo", func() error { _, err := FieldGroupGetInfo(FieldHandle{}); return err }},
		{"field groups", "dcgmFieldGroupGetAll", func() error { _, err := FieldGroupGetAll(); return err }},
		{"topology", "dcgmGetGroupTopology", func() error { _, err := GetGroupTopology(GroupHandle{}); return err }},
		{"architecture", "dcgmGetGpuChipArchitecture", func() error { _, err := GetGPUChipArchitecture(0); return err }},
		{"selection", "dcgmSelectGpusByTopology", func() error { _, err := SelectGPUsByTopology(nil, 1, 0); return err }},
	} {
		t.Run(test.name, func(t *testing.T) {
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, test.call(), &unavailable)
			require.Equal(t, test.symbol, unavailable.Symbol)
		})
	}
}
