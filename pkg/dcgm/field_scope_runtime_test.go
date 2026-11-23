//go:build linux && cgo

package dcgm

import (
	"fmt"
	"os"
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
