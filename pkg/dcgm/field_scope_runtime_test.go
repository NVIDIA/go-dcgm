//go:build linux && cgo

package dcgm

import (
	"fmt"
	"math"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestDCGM48ScopeReads retains entity identity through reads of injected cache values.
func TestDCGM48ScopeReads(t *testing.T) {
	t.Setenv("DCGM_ENABLE_INJECTION_APIS", "1")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	if !dcgmSymbolAvailable("dcgmInjectEntityFieldValue") {
		t.Skip("DCGM runtime lacks entity value injection")
	}
	t.Run("XID", checkXIDScopeReads)
	t.Run("LINK", checkLinkScopeReads)
}

func checkXIDScopeReads(t *testing.T) {
	t.Helper()
	parent := GroupEntityPair{}
	for i, group := range []Field_Entity_Group{FE_GPU, FE_GPU_I, FE_GPU_CI} {
		ids, createErr := CreateFakeEntities([]MigHierarchyInfo{{Entity: GroupEntityPair{EntityGroupId: group}, Parent: parent}})
		require.NoError(t, createErr)
		require.Len(t, ids, 1)
		entity := GroupEntityPair{EntityGroupId: group, EntityId: ids[0]}
		t.Run(group.String(), func(t *testing.T) {
			meta, metaErr := FieldGetByID(DCGM_FI_DEV_XID_ERROR)
			require.NoError(t, metaErr)
			require.True(t, FieldSupportsEntityGroup(DCGM_FI_DEV_XID_ERROR, meta.EntityLevel, group))
			require.NoError(t, injectEntityValueForTest(entity, DCGM_FI_DEV_XID_ERROR, meta.FieldType, int64(70+i)))
			values, readErr := EntitiesGetLatestValues([]GroupEntityPair{entity}, []Short{DCGM_FI_DEV_XID_ERROR}, 0)
			require.NoError(t, readErr)
			require.Len(t, values, 1)
			require.Equal(t, group, values[0].EntityGroupId)
			require.Equal(t, entity.EntityId, values[0].EntityID)
			require.Equal(t, DCGM_ST_OK, values[0].Status)
			require.Equal(t, int64(70+i), values[0].Int64())
		})
		parent = entity
	}
}

func checkLinkScopeReads(t *testing.T) {
	t.Helper()
	if !dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("LINK field metadata 1714 through 1716 requires DCGM 4.8")
	}
	for field := Short(1714); field <= 1716; field++ {
		meta, metaErr := FieldGetByID(field)
		require.NoError(t, metaErr)
		require.Equal(t, FE_LINK, meta.EntityLevel)
	}
	links, linkErr := GetNvLinkLinkStatus()
	require.NoError(t, linkErr)
	var link *NvLinkStatus
	for i := range links {
		if links[i].State == LS_UP {
			link = &links[i]
			break
		}
	}
	if link == nil {
		t.Skip("no active NVLink for LINK field reads")
	}
	// dcgm_link_t packs parent type, 16-bit link index, and 8-bit parent ID.
	entity := GroupEntityPair{EntityGroupId: FE_LINK, EntityId: uint(link.ParentType) | link.Index<<8 | link.ParentId<<24}
	for _, field := range []Short{
		DCGM_FI_PROF_NVLINK_TX_BYTES_PER_SECOND, DCGM_FI_PROF_NVLINK_RX_BYTES_PER_SECOND,
		DCGM_FI_DEV_NVLINK_THROUGHPUT_PER_LINK_BYTES_PER_SECOND, DCGM_FI_DEV_NVLINK_TX_PER_LINK_BYTES_PER_SECOND, DCGM_FI_DEV_NVLINK_RX_PER_LINK_BYTES_PER_SECOND,
	} {
		t.Run(fmt.Sprint(field), func(t *testing.T) {
			meta, metaErr := FieldGetByID(field)
			require.NoError(t, metaErr)
			require.True(t, FieldSupportsEntityGroup(field, meta.EntityLevel, FE_LINK))
			require.NoError(t, injectEntityValueForTest(entity, field, meta.FieldType, 1234))
			values, readErr := EntitiesGetLatestValues([]GroupEntityPair{entity}, []Short{field}, 0)
			require.NoError(t, readErr)
			require.Len(t, values, 1)
			require.Equal(t, FE_LINK, values[0].EntityGroupId)
			require.Equal(t, entity.EntityId, values[0].EntityID)
			require.Equal(t, DCGM_ST_OK, values[0].Status)
			if uint(meta.FieldType) == DCGM_FT_DOUBLE {
				require.Equal(t, float64(1234), values[0].Float64())
			} else {
				require.Equal(t, int64(1234), values[0].Int64())
			}
		})
	}
}

func TestDCGM48InjectionGate(t *testing.T) {
	t.Setenv("DCGM_ENABLE_INJECTION_APIS", "")
	require.NoError(t, os.Unsetenv("DCGM_ENABLE_INJECTION_APIS"))
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	if !dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("injection gate requires DCGM 4.8")
	}
	t.Run("entity injection", func(t *testing.T) {
		injectErr := injectEntityValueForTest(GroupEntityPair{EntityGroupId: FE_GPU}, DCGM_FI_DEV_GPU_TEMP, byte(DCGM_FT_INT64), 40)
		requireDCGMCode(t, injectErr, DCGM_ST_NOT_SUPPORTED)
	})
	_, err = CreateFakeEntities([]MigHierarchyInfo{{Entity: GroupEntityPair{EntityGroupId: FE_GPU}}})
	requireDCGMCode(t, err, DCGM_ST_NOT_SUPPORTED)
}

// TestDeviceMemoryMiB verifies that the legacy attribute layout keeps framebuffer totals in MiB.
func TestDeviceMemoryMiB(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	runOnlyWithLiveGPUs(t)
	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	for _, gpu := range gpus {
		t.Run(fmt.Sprintf("GPU-%d", gpu), func(t *testing.T) {
			device, deviceErr := GetDeviceInfo(gpu)
			require.NoError(t, deviceErr)
			require.NoError(t, WatchFieldValue(gpu, DCGM_FI_DEV_FB_TOTAL, 100*time.Millisecond, time.Second, 1))
			defer func() { require.NoError(t, UnwatchFieldValue(gpu, DCGM_FI_DEV_FB_TOTAL, true)) }()
			require.NoError(t, UpdateAllFields())
			values, readErr := GetLatestValuesForFields(gpu, []Short{DCGM_FI_DEV_FB_TOTAL})
			require.NoError(t, readErr)
			require.Len(t, values, 1)
			require.Equal(t, DCGM_ST_OK, values[0].Status)
			require.Positive(t, values[0].Int64())
			require.EqualValues(t, device.PCI.FBTotal, values[0].Int64())
		})
	}
}

func TestEntityInjectionHelperWithoutInit(t *testing.T) {
	require.Nil(t, dcgmLibHandle)
	err := injectEntityValueForTest(GroupEntityPair{EntityGroupId: FE_GPU}, DCGM_FI_DEV_GPU_TEMP, byte(DCGM_FT_INT64), 40)
	require.True(t, IsFeatureUnavailable(err))
}

// TestLinkFieldMetadataInventory requires the DCGM runtime, but no hostengine or
// GPU. DCGM 4.7 must lack seven reviewed 4.8 additions; newer runtimes must
// provide the full inventory. New native link or vGPU fields require classification.
func TestLinkFieldMetadataInventory(t *testing.T) {
	api := cgoAdapter{}
	if err := api.loadLibrary(); err != nil {
		t.Fatalf("load DCGM runtime: %v", err)
	}
	defer api.closeLibrary()

	if status := FieldsInit(); status != 0 {
		t.Fatalf("initialize field metadata: status %d", status)
	}
	defer func() {
		if status := FieldsTerm(); status != 0 {
			t.Errorf("terminate field metadata: status %d", status)
		}
	}()

	info, err := GetVersionInfo()
	require.NoError(t, err)
	exclusions, err := metadataVersionExclusions(info)
	require.NoError(t, err)
	t.Logf("native metadata library: %s; version exclusions: %v", info.RawBuildInfoString, exclusions)
	t.Run("version_exclusions", func(t *testing.T) {
		for _, fieldID := range exclusions {
			t.Run(fmt.Sprintf("field_%d", fieldID), func(t *testing.T) {
				_, metadataErr := FieldGetByID(fieldID)
				require.Error(t, metadataErr, "reviewed 4.8 field %d must be absent from a 4.7 library", fieldID)
			})
		}
	})

	inventory := linkFieldExpectations()
	for _, fieldID := range dcgmFields {
		if slices.Contains(exclusions, fieldID) {
			continue
		}

		want, listed := inventory[fieldID]
		checkLinkFieldMetadata(t, fieldID, want.declared, want.parent, listed)
	}

	t.Run("native_link_classification", func(t *testing.T) {
		checkNativeLinkClassification(t, inventory)
	})

	t.Run("secondary_alias_metadata", func(t *testing.T) {
		checked := 0
		for _, expected := range boundedAliasExpectations() {
			if slices.Contains(exclusions, expected.field) {
				continue
			}

			checked++
			t.Run(fmt.Sprintf("field_%d", expected.field), func(t *testing.T) {
				meta, err := FieldGetByID(expected.field)
				require.NoError(t, err)
				checkBoundedFieldAlias(t, meta, expected)
			})
		}
		t.Logf("qualified alias fields: %d", checked)
	})
	t.Run("vgpu_metadata", checkVgpuFieldMetadata)
	t.Run("nonlink_metadata", checkNonLinkFieldMetadata)

	t.Logf("checked %d link-capable fields against the loaded DCGM metadata", len(inventory))
}

// checkLinkFieldMetadata verifies one field against the reviewed inventory;
// non-link legacy fields may be absent from the loaded runtime.
func checkLinkFieldMetadata(t *testing.T, fieldID Short, declaredLevel, parentGroup Field_Entity_Group, listed bool) {
	t.Helper()

	meta, err := FieldGetByID(fieldID)
	if err != nil {
		if listed {
			t.Errorf("inventoried link-capable field %d is absent from this runtime: %v", fieldID, err)
		}

		return
	}

	if meta.EntityLevel == FE_LINK && !listed {
		t.Errorf("declared-link field %d (%s) has no parent classification", fieldID, meta.Tag)
	}
	if listed && meta.EntityLevel != declaredLevel {
		t.Errorf("field %d metadata declares %s, inventory declares %s", fieldID, meta.EntityLevel, declaredLevel)
	}

	if !listed {
		return
	}

	for _, group := range []Field_Entity_Group{FE_GPU, FE_SWITCH} {
		parent := GroupEntityPair{EntityGroupId: group}
		entity, conversionErr := LinkEntityPair(0, group, 0)
		require.NoError(t, conversionErr)
		applicable, applicabilityErr := FieldSupportsTarget(meta, entity, parent)
		require.NoError(t, applicabilityErr, "native metadata field %d", fieldID)
		require.Equal(t, group == parentGroup, applicable, "native metadata field %d parent %s", fieldID, group)
	}
}

// checkVgpuFieldMetadata checks all native vGPU declarations against the
// independently reviewed types and complete target/legacy-helper matrix.
func checkVgpuFieldMetadata(t *testing.T) {
	for _, expected := range vgpuFieldExpectations() {
		t.Run(fmt.Sprintf("field_%d", expected.field), func(t *testing.T) {
			meta, err := FieldGetByID(expected.field)
			require.NoError(t, err)
			checkVgpuFieldTargets(t, meta, expected)
		})
	}
	t.Logf("qualified vGPU declarations: %d", len(vgpuFieldExpectations()))
}

// checkNonLinkFieldMetadata verifies actual native metadata for global and GPU
// scopes, independent XID histories, native aliases and a core-only clock.
func checkNonLinkFieldMetadata(t *testing.T) {
	t.Helper()
	tests := []struct {
		name             string
		field            Short
		declared, target Field_Entity_Group
		scope            int
		want             bool
	}{
		{"global on GPU", DCGM_FI_DRIVER_VERSION, FE_NONE, FE_GPU, 0, true},
		{"global on CPU", DCGM_FI_DRIVER_VERSION, FE_NONE, FE_CPU, 0, true},
		{"GPU field at GPU", DCGM_FI_DEV_GPU_TEMP_CELSIUS, FE_GPU, FE_GPU, 1, true},
		{"GPU field at GI", DCGM_FI_DEV_GPU_TEMP_CELSIUS, FE_GPU, FE_GPU_I, 1, true},
		{"GPU field at CI", DCGM_FI_DEV_GPU_TEMP_CELSIUS, FE_GPU, FE_GPU_CI, 1, true},
		{"GPU field at CPU", DCGM_FI_DEV_GPU_TEMP_CELSIUS, FE_GPU, FE_CPU, 1, false},
		{"GI field at GPU", DCGM_FI_PROF_DRAM_UTIL_RATIO, FE_GPU_I, FE_GPU, 1, true},
		{"GI field at GI", DCGM_FI_PROF_DRAM_UTIL_RATIO, FE_GPU_I, FE_GPU_I, 1, true},
		{"GI field at CI", DCGM_FI_PROF_DRAM_UTIL_RATIO, FE_GPU_I, FE_GPU_CI, 1, true},
		{"XID GPU", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU, 1, true},
		{"XID GI", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU_I, 1, true},
		{"XID CI", DCGM_FI_DEV_XID_ERROR, FE_GPU_CI, FE_GPU_CI, 1, true},
		{"clock reasons at CI", DCGM_FI_DEV_CLOCKS_EVENT_REASONS, FE_GPU, FE_GPU_CI, 1, true},
		{"core clock at core", DCGM_FI_DEV_CPU_CLOCK_CURRENT, FE_CPU_CORE, FE_CPU_CORE, 1, true},
		{"core clock at CPU", DCGM_FI_DEV_CPU_CLOCK_CURRENT, FE_CPU_CORE, FE_CPU, 1, false},
	}
	t.Logf("qualified non-link cases: %d", len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := FieldGetByID(tt.field)
			require.NoError(t, err)
			require.Equal(t, tt.field, meta.FieldID)
			require.Equal(t, tt.declared, meta.EntityLevel)
			require.Equal(t, tt.scope, meta.Scope)
			got, err := FieldSupportsTarget(meta, GroupEntityPair{EntityGroupId: tt.target}, GroupEntityPair{})
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// checkNativeLinkClassification scans native field IDs, including IDs not in the
// generated catalog, so new link or vGPU declarations cannot silently evade the gate.
func checkNativeLinkClassification(t *testing.T, inventory map[Short]struct{ declared, parent Field_Entity_Group }) {
	t.Helper()
	links, vgpus := 0, 0
	vgpuInventory := make(map[Short]bool)
	for _, expected := range vgpuFieldExpectations() {
		vgpuInventory[expected.field] = true
	}
	for fieldID := uint32(1); fieldID <= math.MaxUint16; fieldID++ {
		meta, err := FieldGetByID(Short(fieldID))
		if err != nil {
			continue
		}

		switch meta.EntityLevel {
		case FE_LINK:
			_, classified := inventory[meta.FieldID]
			require.True(t, classified, "native declared-link field %d (%s) has no parent classification", meta.FieldID, meta.Tag)
			links++
		case FE_VGPU:
			require.True(t, vgpuInventory[meta.FieldID], "native vGPU field %d (%s) has no target classification", meta.FieldID, meta.Tag)
			vgpus++
		}
	}

	t.Logf("classified native declared-link fields: %d; vGPU declarations: %d", links, vgpus)
	require.Equal(t, len(vgpuInventory), vgpus)
}
