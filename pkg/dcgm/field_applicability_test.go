package dcgm

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFieldSupportsTarget checks independent metadata/target expectations from
// DCGM 4.7 CacheManager practical-entity dispatch and sysmon field readers.
func TestFieldSupportsTarget(t *testing.T) {
	tests := []struct {
		name   string
		field  Short
		level  Field_Entity_Group
		target Field_Entity_Group
		want   bool
	}{
		{"GPU temperature", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_GPU, true},
		{"GPU field on GI aliases physical GPU", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_GPU_I, true},
		{"GPU field on CI aliases physical GPU", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_GPU_CI, true},
		{"GI field on GPU", DCGM_FI_PROF_DRAM_UTIL_RATIO, FE_GPU_I, FE_GPU, true},
		{"GI field on CI aliases GI", DCGM_FI_PROF_DRAM_UTIL_RATIO, FE_GPU_I, FE_GPU_CI, true},
		{"XID physical history", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU, true},
		{"XID GI history", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU_I, true},
		{"XID CI history", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU_CI, true},
		{"clock reasons native CI alias", DCGM_FI_DEV_CLOCKS_EVENT_REASONS, FE_GPU, FE_GPU_CI, true},
		{"GPU field on CPU", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_CPU, false},
		{"CPU socket temperature", DCGM_FI_DEV_CPU_TEMP_CURRENT, FE_CPU, FE_CPU, true},
		{"CPU temperature on core", DCGM_FI_DEV_CPU_TEMP_CELSIUS, FE_CPU, FE_CPU_CORE, true},
		{"CPU power on core", DCGM_FI_DEV_CPU_POWER_WATTS, FE_CPU, FE_CPU_CORE, true},
		{"CPU vendor on core", DCGM_FI_DEV_CPU_VENDOR, FE_CPU, FE_CPU_CORE, true},
		{"core utilization", DCGM_FI_DEV_CPU_UTIL_TOTAL, FE_CPU_CORE, FE_CPU_CORE, true},
		{"CPU aggregate utilization", DCGM_FI_DEV_CPU_UTIL_TOTAL, FE_CPU_CORE, FE_CPU, true},
		{"core clock", DCGM_FI_DEV_CPU_CLOCK_CURRENT, FE_CPU_CORE, FE_CPU_CORE, true},
		{"CPU aggregate clock unsupported", DCGM_FI_DEV_CPU_CLOCK_CURRENT, FE_CPU_CORE, FE_CPU, false},
		{"switch port throughput as aggregate", DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX, FE_LINK, FE_SWITCH, true},
		{"switch port bytes as aggregate", DCGM_FI_DEV_NVSWITCH_LINK_TX_BYTES_TOTAL, FE_LINK, FE_SWITCH, true},
		{"switch port state not aggregate", DCGM_FI_DEV_NVSWITCH_LINK_STATUS, FE_LINK, FE_SWITCH, false},
		{"GPU link state not switch aggregate", DCGM_FI_DEV_NVLINK_LINK_STATUS, FE_LINK, FE_SWITCH, false},
		{"switch aggregate", DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX, FE_SWITCH, FE_SWITCH, true},
		{"CPU field on GPU", DCGM_FI_DEV_CPU_UTIL_TOTAL, FE_CPU_CORE, FE_GPU, false},
		{"future nonlink field", 65000, FE_GPU, FE_GPU, true},
		{"future nonlink field incompatible target", 65000, FE_SWITCH, FE_CPU, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := FieldMeta{FieldID: tt.field, Scope: 1, EntityLevel: tt.level}
			entity := GroupEntityPair{EntityGroupId: tt.target, EntityId: 0}
			got, err := FieldSupportsTarget(meta, entity, GroupEntityPair{})
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestFieldSupportsTargetLinks checks the independent existing full inventory
// through the precise boundary and native packed GPU/switch targets.
func TestFieldSupportsTargetLinks(t *testing.T) {
	for fieldID, expected := range linkFieldExpectations() {
		for _, parentGroup := range []Field_Entity_Group{FE_GPU, FE_SWITCH} {
			parent := GroupEntityPair{EntityGroupId: parentGroup, EntityId: 0}
			entity, err := LinkEntityPair(17, parentGroup, parent.EntityId)
			require.NoError(t, err)
			got, err := FieldSupportsTarget(FieldMeta{FieldID: fieldID, Scope: 1, EntityLevel: expected.declared}, entity, parent)
			require.NoError(t, err)
			require.Equal(t, expected.parent == parentGroup, got, "field %d parent %s", fieldID, parentGroup)
		}
	}
}

// TestFieldSupportsTargetGlobals checks globals on real targets including both
// link parents; FE_NONE never substitutes for a missing link parent.
func TestFieldSupportsTargetGlobals(t *testing.T) {
	meta := FieldMeta{FieldID: DCGM_FI_DRIVER_VERSION, Scope: 0, EntityLevel: FE_NONE}
	for group := FE_NONE; group < FE_COUNT; group++ {
		if group == FE_LINK {
			continue
		}

		got, err := FieldSupportsTarget(meta, GroupEntityPair{EntityGroupId: group}, GroupEntityPair{})
		require.NoError(t, err)
		require.True(t, got)
	}
	for _, group := range []Field_Entity_Group{FE_GPU, FE_SWITCH} {
		entity, err := LinkEntityPair(0, group, 0)
		require.NoError(t, err)
		got, err := FieldSupportsTarget(meta, entity, GroupEntityPair{EntityGroupId: group})
		require.NoError(t, err)
		require.True(t, got)
	}
}

// TestFieldSupportsTargetRejectsInvalid distinguishes invalid/unknown inputs
// from a valid target that cannot collect a known field.
func TestFieldSupportsTargetRejectsInvalid(t *testing.T) {
	gpuMeta := FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 1, EntityLevel: FE_GPU}
	gpu := GroupEntityPair{EntityGroupId: FE_GPU}
	link := GroupEntityPair{EntityGroupId: FE_LINK, EntityId: 1}
	tests := []struct {
		name   string
		meta   FieldMeta
		entity GroupEntityPair
		parent GroupEntityPair
	}{
		{"missing metadata", FieldMeta{}, gpu, GroupEntityPair{}},
		{"unknown scope", FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 2, EntityLevel: FE_GPU}, gpu, GroupEntityPair{}},
		{"global with entity level", FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 0, EntityLevel: FE_GPU}, gpu, GroupEntityPair{}},
		{"entity without level", FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 1, EntityLevel: FE_NONE}, gpu, GroupEntityPair{}},
		{"unknown declared level", FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 1, EntityLevel: FE_COUNT}, gpu, GroupEntityPair{}},
		{"unknown target", gpuMeta, GroupEntityPair{EntityGroupId: FE_COUNT}, GroupEntityPair{}},
		{"unknown parent", gpuMeta, gpu, GroupEntityPair{EntityGroupId: FE_COUNT}},
		{"missing link parent", gpuMeta, link, GroupEntityPair{}},
		{"wrong packed parent", gpuMeta, link, GroupEntityPair{EntityGroupId: FE_SWITCH}},
		{"wrong packed parent ID", gpuMeta, link, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 1}},
		{"unclassified declared link", FieldMeta{FieldID: 65000, Scope: 1, EntityLevel: FE_LINK}, link, gpu},
		{"unclassified new GPU link candidate", FieldMeta{FieldID: 65000, Scope: 1, EntityLevel: FE_GPU}, link, gpu},
		{"missing parent group with ID", gpuMeta, gpu, GroupEntityPair{EntityId: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FieldSupportsTarget(tt.meta, tt.entity, tt.parent)
			require.Error(t, err)
			require.False(t, got)
		})
	}
	if uint64(math.MaxUint) > math.MaxUint32 {
		overflow := uint64(math.MaxUint32) + 1
		got, err := FieldSupportsTarget(gpuMeta, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: uint(overflow)}, GroupEntityPair{})
		require.Error(t, err)
		require.False(t, got)
	}

	got, err := FieldSupportsTarget(gpuMeta, link, gpu)
	require.NoError(t, err)
	require.False(t, got, "temperature is known and incompatible with a GPU link")
}

// TestNativeEntityKeysStayDistinct checks applicability using equal integer IDs
// in separate native groups, without converting a native ID to an NVML index.
func TestNativeEntityKeysStayDistinct(t *testing.T) {
	meta := FieldMeta{FieldID: DCGM_FI_DEV_GPU_TEMP, Scope: 1, EntityLevel: FE_GPU}
	for _, id := range []uint{0, 1, math.MaxUint32} {
		for _, group := range []Field_Entity_Group{FE_GPU, FE_GPU_I, FE_GPU_CI, FE_SWITCH, FE_CPU, FE_CPU_CORE} {
			entity := GroupEntityPair{EntityGroupId: group, EntityId: id}
			got, err := FieldSupportsTarget(meta, entity, GroupEntityPair{})
			require.NoError(t, err)
			require.Equal(t, group == FE_GPU || group == FE_GPU_I || group == FE_GPU_CI, got)
			require.Equal(t, id, entity.EntityId)
		}
	}
}

// TestLegacyHelpersKeepDeclaredScopes ensures the precise native aliases do
// not broaden the older helpers' documented declared-scope contracts.
func TestLegacyHelpersKeepDeclaredScopes(t *testing.T) {
	for _, tt := range []struct {
		field         Short
		level, target Field_Entity_Group
	}{
		{DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_GPU_CI},
		{DCGM_FI_DEV_CPU_UTIL_RATIO, FE_CPU_CORE, FE_CPU},
		{DCGM_FI_DEV_CPU_TEMP_CELSIUS, FE_CPU, FE_CPU_CORE},
		{DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX, FE_LINK, FE_SWITCH},
	} {
		require.False(t, FieldSupportsEntityGroup(tt.field, tt.level, tt.target))
		require.False(t, FieldSupportsEntityGroupWithParent(tt.field, tt.level, tt.target, FE_NONE))
	}
	require.True(t, FieldSupportsEntityGroup(DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU))
	require.True(t, FieldSupportsEntityGroupWithParent(DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU_I, FE_NONE))
}

// boundedAliasExpectation records independently reviewed native metadata and
// secondary targets. All link targets in these cases have switch parents.
type boundedAliasExpectation struct {
	field     Short
	declared  Field_Entity_Group
	secondary []Field_Entity_Group
}

// boundedAliasExpectations lists every bounded exception from DCGM 4.8
// b6c33f98 Sysmon CalculateUtilizationForEntity/UpdateFields and NVSDM
// UpdateFieldsFromNvswitchLibrary/getNvsdmPortFieldId/isCompositeFieldId.
func boundedAliasExpectations() []boundedAliasExpectation {
	return []boundedAliasExpectation{
		{DCGM_FI_DEV_CPU_UTIL_RATIO, FE_CPU_CORE, []Field_Entity_Group{FE_CPU}},
		{DCGM_FI_DEV_CPU_UTIL_USER_RATIO, FE_CPU_CORE, []Field_Entity_Group{FE_CPU}},
		{DCGM_FI_DEV_CPU_UTIL_NICE_RATIO, FE_CPU_CORE, []Field_Entity_Group{FE_CPU}},
		{DCGM_FI_DEV_CPU_UTIL_SYS_RATIO, FE_CPU_CORE, []Field_Entity_Group{FE_CPU}},
		{DCGM_FI_DEV_CPU_UTIL_IRQ_RATIO, FE_CPU_CORE, []Field_Entity_Group{FE_CPU}},
		{DCGM_FI_DEV_CPU_TEMP_CELSIUS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_TEMP_WARNING_CELSIUS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_TEMP_CRITICAL_CELSIUS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_POWER_WATTS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_SYSIO_POWER_WATTS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_MODULE_POWER_WATTS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_POWER_LIMIT_WATTS, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_VENDOR, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_CPU_MODEL, FE_CPU, []Field_Entity_Group{FE_CPU_CORE}},
		{DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_LINK_TX_BYTES_TOTAL, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_LINK_RX_BYTES_TOTAL, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERROR_TOTAL, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERROR_TOTAL, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_UUID, FE_LINK, []Field_Entity_Group{FE_SWITCH}},
		{DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX, FE_SWITCH, []Field_Entity_Group{FE_LINK}},
		{DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX, FE_SWITCH, []Field_Entity_Group{FE_LINK}},
		{DCGM_FI_DEV_NVSWITCH_TX_BYTES_TOTAL, FE_SWITCH, []Field_Entity_Group{FE_LINK}},
		{DCGM_FI_DEV_NVSWITCH_RX_BYTES_TOTAL, FE_SWITCH, []Field_Entity_Group{FE_LINK}},
	}
}

// TestBoundedFieldAliases checks every bounded exception on every known native
// group, both link parents and the two legacy public helpers.
func TestBoundedFieldAliases(t *testing.T) {
	for _, expected := range boundedAliasExpectations() {
		t.Run(fmt.Sprintf("field_%d", expected.field), func(t *testing.T) {
			meta := FieldMeta{FieldID: expected.field, Scope: 1, EntityLevel: expected.declared}
			checkBoundedFieldAlias(t, meta, expected)
		})
	}
}

// checkBoundedFieldAlias applies explicit fixture expectations to all native
// target groups; the fixture is shared with the library-only metadata gate.
func checkBoundedFieldAlias(t *testing.T, meta FieldMeta, expected boundedAliasExpectation) {
	t.Helper()
	require.Equal(t, expected.field, meta.FieldID)
	require.Equal(t, expected.declared, meta.EntityLevel)
	require.Equal(t, 1, meta.Scope)
	for group := FE_NONE; group < FE_COUNT; group++ {
		parents := []Field_Entity_Group{FE_NONE}
		if group == FE_LINK {
			parents = []Field_Entity_Group{FE_GPU, FE_SWITCH}
		}

		for _, parent := range parents {
			t.Run(fmt.Sprintf("target_%d/parent_%d", group, parent), func(t *testing.T) {
				entity := GroupEntityPair{EntityGroupId: group}
				if group == FE_LINK {
					var err error
					entity, err = LinkEntityPair(0, parent, 0)
					require.NoError(t, err)
				}

				got, err := FieldSupportsTarget(meta, entity, GroupEntityPair{EntityGroupId: parent})
				require.NoError(t, err)
				primary := group == expected.declared
				secondary := slices.Contains(expected.secondary, group)
				validParent := group != FE_LINK || parent == FE_SWITCH
				require.Equal(t, (primary || secondary) && validParent, got)
				require.Equal(t, primary, FieldSupportsEntityGroup(expected.field, expected.declared, group))
				require.Equal(t, primary && validParent, FieldSupportsEntityGroupWithParent(expected.field, expected.declared, group, parent))
			})
		}
	}
}

// metadataVersionExclusions selects only the reviewed 4.7 metadata exclusions.
// Newer runtimes must supply the complete inventory; older versions fail.
func metadataVersionExclusions(info VersionInfo) ([]Short, error) {
	version, err := metadataRuntimeVersion(info)
	if err != nil {
		return nil, err
	}

	if version[0] == 4 && version[1] == 7 {
		return []Short{
			DCGM_FI_DEV_NVSWITCH_LINK_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_LINK_RX_BYTES_TOTAL,
			DCGM_FI_DEV_NVSWITCH_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_RX_BYTES_TOTAL,
			DCGM_FI_DEV_NVLINK_THROUGHPUT_PER_LINK_BYTES_PER_SECOND,
			DCGM_FI_DEV_NVLINK_TX_PER_LINK_BYTES_PER_SECOND,
			DCGM_FI_DEV_NVLINK_RX_PER_LINK_BYTES_PER_SECOND,
		}, nil
	}

	if version[0] > 4 || (version[0] == 4 && version[1] >= 8) {
		return nil, nil
	}

	return nil, fmt.Errorf("metadata qualification requires DCGM 4.7 or newer, got %d.%d.%d", version[0], version[1], version[2])
}

// metadataRuntimeVersion parses the library's version key without guessing from
// another build-info value, accepting exactly three unsigned numeric components.
func metadataRuntimeVersion(info VersionInfo) ([3]uint64, error) {
	var version [3]uint64
	var text string
	for _, pair := range strings.Split(info.RawBuildInfoString, ";") {
		key, value, found := strings.Cut(pair, ":")
		if !found || strings.TrimSpace(key) != "version" {
			continue
		}

		if text != "" || strings.TrimSpace(value) == "" {
			return version, errors.New("invalid or duplicate native version key")
		}

		text = strings.TrimSpace(value)
	}

	parts := strings.Split(text, ".")
	if len(parts) != len(version) {
		return version, fmt.Errorf("native version must have three numeric components: %q", text)
	}

	for i, part := range parts {
		value, err := strconv.ParseUint(part, 10, 32)
		if err != nil || (len(part) > 1 && part[0] == '0') {
			return version, fmt.Errorf("invalid native version component %q in %q", part, text)
		}

		version[i] = value
	}

	return version, nil
}

// TestMetadataVersionExclusions checks exact baseline selection and rejects
// malformed, ambiguous and unsupported native versions without a fallback.
func TestMetadataVersionExclusions(t *testing.T) {
	excluded := []Short{
		DCGM_FI_DEV_NVSWITCH_LINK_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_LINK_RX_BYTES_TOTAL,
		DCGM_FI_DEV_NVSWITCH_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_RX_BYTES_TOTAL,
		DCGM_FI_DEV_NVLINK_THROUGHPUT_PER_LINK_BYTES_PER_SECOND,
		DCGM_FI_DEV_NVLINK_TX_PER_LINK_BYTES_PER_SECOND,
		DCGM_FI_DEV_NVLINK_RX_PER_LINK_BYTES_PER_SECOND,
	}
	tests := []struct {
		name, raw string
		want      []Short
		invalid   bool
	}{
		{"4.7.0", "version:4.7.0;arch:x86_64", excluded, false},
		{"CI 4.7.1", "arch:x86_64; version : 4.7.1 ;buildid:123", excluded, false},
		{"4.7 patch", "version:4.7.99", excluded, false},
		{"4.8", "version:4.8.0", nil, false},
		{"newer minor", "version:4.99.1", nil, false},
		{"newer major", "version:5.0.0", nil, false},
		{"older major", "version:3.99.0", nil, true},
		{"older minor", "version:4.6.9", nil, true},
		{"missing", "arch:4.7.0", nil, true},
		{"empty", "version:", nil, true},
		{"duplicate", "version:4.7.0;version:4.8.0", nil, true},
		{"missing patch", "version:4.7", nil, true},
		{"non-numeric", "version:4.x.0", nil, true},
		{"signed", "version:+4.7.0", nil, true},
		{"negative", "version:-4.7.0", nil, true},
		{"leading zero", "version:4.07.0", nil, true},
		{"suffix", "version:4.7.0-rc", nil, true},
		{"overflow", "version:4294967296.7.0", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := metadataVersionExclusions(VersionInfo{RawBuildInfoString: tt.raw})
			if tt.invalid {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// vgpuFieldExpectation records each vGPU declaration and its independently
// verified GPU-updater dispatch in DCGM 4.7 63d15148 and 4.8 b6c33f98.
type vgpuFieldExpectation struct {
	field      Short
	fieldType  byte
	gpuTargets bool
}

// vgpuFieldExpectations covers every FE_VGPU declaration in the pinned native
// sources. Instance-only fields have no GPU/GI/CI updater implementation.
func vgpuFieldExpectations() []vgpuFieldExpectation {
	return []vgpuFieldExpectation{
		{DCGM_FI_DEV_VGPU_SUPPORTED_INFO, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_CREATABLE_IDS, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_INSTANCE_INFO, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_UTIL_INFO, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_PROCESS_UTIL_INFO, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_SUPPORTED_IDS, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_TYPE_INFO, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_TYPE_NAME, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_TYPE_CLASS, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_TYPE_LICENSE, byte(DCGM_FT_BINARY), true},
		{DCGM_FI_DEV_VGPU_LICENSE_STATUS, byte(DCGM_FT_INT64), true},
		{DCGM_FI_DEV_VGPU_VM_ID, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_VM_NAME, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_TYPE, byte(DCGM_FT_INT64), false},
		{DCGM_FI_DEV_VGPU_UUID, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_DRIVER_VERSION, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_MEMORY_USAGE, byte(DCGM_FT_INT64), false},
		{DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT, byte(DCGM_FT_INT64), false},
		{DCGM_FI_DEV_VGPU_ENC_STATS, byte(DCGM_FT_BINARY), false},
		{DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO, byte(DCGM_FT_BINARY), false},
		{DCGM_FI_DEV_VGPU_FBC_STATS, byte(DCGM_FT_BINARY), false},
		{DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO, byte(DCGM_FT_BINARY), false},
		{DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATUS, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_PCI_ID, byte(DCGM_FT_STRING), false},
		{DCGM_FI_DEV_VGPU_GPU_INSTANCE_ID, byte(DCGM_FT_INT64), false},
		{DCGM_FI_DEV_VGPU_FB_USED_BYTES, byte(DCGM_FT_INT64), false},
	}
}

// TestVgpuFieldTargets checks all vGPU declarations on every native target,
// including the license-status scalar and binary fields used as metadata.
func TestVgpuFieldTargets(t *testing.T) {
	for _, expected := range vgpuFieldExpectations() {
		t.Run(fmt.Sprintf("field_%d", expected.field), func(t *testing.T) {
			meta := FieldMeta{FieldID: expected.field, FieldType: expected.fieldType, Scope: 1, EntityLevel: FE_VGPU}
			checkVgpuFieldTargets(t, meta, expected)
		})
	}
}

// checkVgpuFieldTargets verifies native metadata and reuses the exhaustive
// target and legacy-helper matrix without changing exporter reporting policy.
func checkVgpuFieldTargets(t *testing.T, meta FieldMeta, expected vgpuFieldExpectation) {
	t.Helper()
	require.Equal(t, expected.fieldType, meta.FieldType)
	alias := boundedAliasExpectation{field: expected.field, declared: FE_VGPU}
	if expected.gpuTargets {
		alias.secondary = []Field_Entity_Group{FE_GPU, FE_GPU_I, FE_GPU_CI}
	}

	checkBoundedFieldAlias(t, meta, alias)
}
