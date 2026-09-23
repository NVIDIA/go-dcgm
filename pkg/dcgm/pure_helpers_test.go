package dcgm

import (
	"strings"
	"testing"
)

func TestSmallMappingHelpers(t *testing.T) {
	for location, want := range map[int]string{0: "L1", 1: "L2", 2: "Device", 3: "Register", 4: "Texture", 5: "N/A"} {
		if got := dbeLocation(location); got != want {
			t.Errorf("dbeLocation(%d) = %q, want %q", location, got, want)
		}
	}
	for link, want := range map[P2PLinkType]struct {
		count uint
		ok    bool
	}{
		P2PLinkUnknown:        {0, false},
		SingleNVLINKLink:      {1, true},
		SingleNVLINKLink + 35: {36, true},
		SingleNVLINKLink + P2PLinkType(maxNVLinkCount): {0, false},
	} {
		gotCount, gotOK := link.nvLinkCount()
		if gotCount != want.count || gotOK != want.ok {
			t.Errorf("nvLinkCount(%d) = (%d, %v), want (%d, %v)", link, gotCount, gotOK, want.count, want.ok)
		}
	}
}

func TestFieldValuePureHelpers(t *testing.T) {
	var ascii [4096]byte
	for i := range ascii {
		ascii[i] = 'A'
	}
	if got := FindFirstNonAsciiIndex(ascii); got != len(ascii) {
		t.Fatalf("all-ASCII index = %d, want %d", got, len(ascii))
	}
	ascii[17] = 0
	if got := FindFirstNonAsciiIndex(ascii); got != 17 {
		t.Fatalf("control-byte index = %d, want 17", got)
	}
	ascii[17] = 'A'
	ascii[23] = 0xff
	if got := FindFirstNonAsciiIndex(ascii); got != 23 {
		t.Fatalf("high-byte index = %d, want 23", got)
	}

	text := "field text"
	stringValue := FieldValue_v2{FieldType: DCGM_FT_STRING, StringValue: &text}
	if got := Fv2_String(stringValue); got != text {
		t.Fatalf("Fv2_String(string) = %q", got)
	}
	var raw FieldValue_v2
	copy(raw.Value[:], "raw")
	if got := Fv2_String(raw); !strings.HasPrefix(got, "raw") || len(got) != len(raw.Value) {
		t.Fatalf("Fv2_String(raw) has length %d and prefix %q", len(got), got[:3])
	}
	if got := Fv2_Blob(raw); got != raw.Value {
		t.Fatal("Fv2_Blob changed the field bytes")
	}

	if got := ToFieldMeta(nil); got != (FieldMeta{}) {
		t.Fatalf("ToFieldMeta(nil) = %+v", got)
	}
}

func TestCreateFakeEntitiesBounds(t *testing.T) {
	ids, err := CreateFakeEntities(nil)
	if err != nil || len(ids) != 0 {
		t.Fatalf("CreateFakeEntities(nil) = (%v, %v)", ids, err)
	}

	entities := make([]MigHierarchyInfo, maxFakeEntities+1)
	for i := range entities {
		entities[i].Entity.EntityId = uint(i)
	}
	got := boundedFakeEntities(entities)
	if len(got) != maxFakeEntities {
		t.Fatalf("bounded entity count = %d, want %d", len(got), maxFakeEntities)
	}
	if got[len(got)-1].Entity.EntityId != uint(len(got)-1) {
		t.Fatal("bounding discarded or reordered an entity within the limit")
	}
}

func TestInjectFieldValueRejectsInvalidInputBeforeDCGMCall(t *testing.T) {
	for _, test := range []struct {
		name      string
		fieldType uint
		value     any
		want      string
	}{
		{name: "wrong integer type", fieldType: DCGM_FT_INT64, value: "invalid", want: "requires int64"},
		{name: "wrong float type", fieldType: DCGM_FT_DOUBLE, value: int64(1), want: "requires float64"},
		{name: "unsupported field type", fieldType: 0, value: int64(1), want: "unsupported field type"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := InjectFieldValue(0, 0, test.fieldType, 0, 0, test.value); err == nil ||
				!strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}
