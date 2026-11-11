//go:build linux && cgo

package dcgm

import "testing"

// TestLinkFieldMetadataInventory requires the DCGM runtime, but no hostengine or
// GPU. It can run against DCGM 4.7 or 4.8; every declared-link field known to
// this binding must be classified in the independently reviewed inventory.
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

	inventory := linkFieldExpectations()
	for _, fieldID := range dcgmFields {
		want, listed := inventory[fieldID]
		checkLinkFieldMetadata(t, fieldID, want.declared, listed)
	}

	t.Logf("checked %d link-capable fields against the loaded DCGM metadata", len(inventory))
}

// checkLinkFieldMetadata verifies one field against the reviewed inventory;
// non-link legacy fields may be absent from the loaded runtime.
func checkLinkFieldMetadata(t *testing.T, fieldID Short, declaredLevel Field_Entity_Group, listed bool) {
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
}
