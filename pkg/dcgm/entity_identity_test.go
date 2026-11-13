package dcgm

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLinkEntityPair checks the dcgm_link_t representation independently of
// either adapter, including native ID zero and both parent types.
func TestLinkEntityPair(t *testing.T) {
	tests := []struct {
		name    string
		index   uint
		parent  Field_Entity_Group
		id      uint
		want    uint
		invalid bool
	}{
		{"GPU zero", 0, FE_GPU, 0, 0x00000001, false},
		{"switch zero", 0, FE_SWITCH, 0, 0x00000003, false},
		{"sparse GPU port", 17, FE_GPU, 3, 0x03001101, false},
		{"bounds", 65535, FE_SWITCH, 255, 0xffffff03, false},
		{"index overflow", 65536, FE_GPU, 0, 0, true},
		{"parent overflow", 0, FE_GPU, 256, 0, true},
		{"unknown parent", 0, FE_COUNT, 0, 0, true},
		{"missing parent", 0, FE_NONE, 0, 0, true},
		{"wrong parent", 0, FE_GPU_CI, 0, 0, true},
		{"wide index", math.MaxUint, FE_GPU, 0, 0, true},
		{"wide parent", 0, FE_GPU, math.MaxUint, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LinkEntityPair(tt.index, tt.parent, tt.id)
			if tt.invalid {
				require.Error(t, err)
				require.Zero(t, got)
				var nativeErr *Error
				require.ErrorAs(t, err, &nativeErr)
				require.EqualValues(t, DCGM_ST_BADPARAM, nativeErr.Code)
				return
			}

			require.NoError(t, err)
			require.Equal(t, GroupEntityPair{EntityGroupId: FE_LINK, EntityId: tt.want}, got)
		})
	}
}

// TestLinkIdentityAdapters compares independently expected keys at group and
// read boundaries and verifies that returned errors and samples are untouched.
func TestLinkIdentityAdapters(t *testing.T) {
	for _, parent := range []Field_Entity_Group{FE_GPU, FE_SWITCH} {
		for _, index := range []uint{0, 17, 65535} {
			want := GroupEntityPair{EntityGroupId: FE_LINK, EntityId: uint(parent) | index*256 | 255*16777216}
			var added, read GroupEntityPair
			sentinel := errors.New("native failure")
			err := addLinkEntityToGroup(GroupHandle{}, index, parent, 255,
				func(_ GroupHandle, group Field_Entity_Group, id uint) error {
					added = GroupEntityPair{EntityGroupId: group, EntityId: id}
					return sentinel
				})
			require.ErrorIs(t, err, sentinel)
			fields := []Short{DCGM_FI_DEV_NVLINK_LINK_STATUS}
			samples := []FieldValue_v1{{FieldID: fields[0], TS: 12345, Status: 7}}
			got, err := linkGetLatestValues(index, parent, 255, fields,
				func(group Field_Entity_Group, id uint, requested []Short) ([]FieldValue_v1, error) {
					read = GroupEntityPair{EntityGroupId: group, EntityId: id}
					require.Equal(t, fields, requested)
					return samples, sentinel
				})
			require.ErrorIs(t, err, sentinel)
			require.Equal(t, samples, got)
			require.Equal(t, want, added)
			require.Equal(t, want, read)
		}
	}
}

// TestInvalidLinkAdapters proves invalid identity input fails before native calls.
func TestInvalidLinkAdapters(t *testing.T) {
	for _, input := range []NvLinkStatus{
		{Index: 65536, ParentType: FE_GPU},
		{ParentId: 256, ParentType: FE_SWITCH},
		{ParentType: FE_CPU},
	} {
		err := addLinkEntityToGroup(GroupHandle{}, input.Index, input.ParentType, input.ParentId,
			func(GroupHandle, Field_Entity_Group, uint) error { t.Fatal("native admission called"); return nil })
		require.Error(t, err)
		_, err = linkGetLatestValues(input.Index, input.ParentType, input.ParentId, []Short{DCGM_FI_DEV_GPU_TEMP},
			func(Field_Entity_Group, uint, []Short) ([]FieldValue_v1, error) {
				t.Fatal("native read called")
				return nil, nil
			})
		require.Error(t, err)
		// Public adapters must also stop before entering an uninitialized DCGM library.
		require.Error(t, AddLinkEntityToGroup(GroupHandle{}, input.Index, input.ParentType, input.ParentId))
		_, err = LinkGetLatestValues(input.Index, input.ParentType, input.ParentId, []Short{DCGM_FI_DEV_GPU_TEMP})
		require.Error(t, err)
	}
}
