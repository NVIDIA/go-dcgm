//go:build linux && cgo

package dcgm

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLinkFieldMetadataInventory checks link and vGPU fields in the loaded
// DCGM library. It expects all known fields except the seven that are absent
// in DCGM 4.7. This test needs no running host engine or GPU.
func TestLinkFieldMetadataInventory(t *testing.T) {
	t.Run("terminated_catalog", testFieldMetadataAfterFieldsTerm)

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
	t.Logf("DCGM library: %s; fields absent in this version: %v", info.RawBuildInfoString, exclusions)

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
		checkLinkFieldClassification(t, inventory)
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

// checkLinkFieldMetadata compares a field with the link inventory.
// Older libraries may omit fields outside that inventory.
func checkLinkFieldMetadata(t *testing.T, fieldID Short, declaredLevel, parentGroup Field_Entity_Group, listed bool) {
	t.Helper()

	meta, err := FieldGetByID(fieldID)
	if err != nil {
		if listed {
			t.Errorf("DCGM library is missing field %d listed in the link inventory: %v", fieldID, err)
		}

		return
	}

	if meta.EntityLevel == FE_LINK && !listed {
		t.Errorf("field %d (%s) has link scope but no GPU or switch parent in the inventory", fieldID, meta.Tag)
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
		require.NoError(t, applicabilityErr, "DCGM field %d", fieldID)
		require.Equal(t, group == parentGroup, applicable, "DCGM field %d parent %s", fieldID, group)
	}
}

// checkVgpuFieldMetadata checks the scope of each vGPU field and which devices it can target.
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

// checkNonLinkFieldMetadata checks which devices can use each non-link field in the table.
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

// checkLinkFieldClassification checks every field ID the library can report.
// A new link or vGPU field must be added to the inventory.
func checkLinkFieldClassification(t *testing.T, inventory map[Short]struct{ declared, parent Field_Entity_Group }) {
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
			require.True(t, classified, "DCGM link field %d (%s) has no GPU or switch parent in the inventory", meta.FieldID, meta.Tag)
			links++
		case FE_VGPU:
			require.True(t, vgpuInventory[meta.FieldID], "DCGM vGPU field %d (%s) has no entry in the inventory", meta.FieldID, meta.Tag)
			vgpus++
		}
	}

	t.Logf("checked %d DCGM link fields and %d vGPU fields", links, vgpus)
	require.Equal(t, len(vgpuInventory), vgpus)
}

// testFieldMetadataAfterFieldsTerm checks how an unknown field is reported
// before and after restarting the catalog within a DCGM session.
func testFieldMetadataAfterFieldsTerm(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	_, err = FieldGetByID(DCGM_FI_DEV_GPU_TEMP)
	require.NoError(t, err)

	require.Zero(t, FieldsTerm())
	// Restore the catalog before shutdown even if an assertion stops the test.
	defer func() { require.Zero(t, FieldsInit()) }()
	require.Positive(t, dcgmInitCounter)

	metadata, err := FieldGetByID(65535)
	require.Empty(t, metadata)
	require.Error(t, err)
	require.NotErrorIs(t, err, errors.ErrUnsupported)

	require.Zero(t, FieldsInit())
	_, err = FieldGetByID(65535)
	require.ErrorIs(t, err, errors.ErrUnsupported)
}
