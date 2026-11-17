//go:build linux && cgo

package dcgm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGroupAdmissionPreservesNativeStatus loads DCGM without initializing it,
// so zero handles must preserve the real native bad-parameter error.
func TestGroupAdmissionPreservesNativeStatus(t *testing.T) {
	// Use zero handles even if an earlier runtime test left its session handle.
	previousHandle := handle
	handle = dcgmHandle{}
	defer func() { handle = previousHandle }()

	api := cgoAdapter{}
	require.NoError(t, api.loadLibrary(), "load DCGM runtime")
	defer api.closeLibrary()

	tests := []struct {
		name string
		add  func() error
	}{
		{"entity", func() error { return AddEntityToGroup(GroupHandle{}, FE_LINK, uint(FE_GPU)) }},
		{"link", func() error { return AddLinkEntityToGroup(GroupHandle{}, 0, FE_GPU, 0) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.add()
			wantMessage := fmt.Sprintf("error adding entity group type %v, entity %v to group: %s",
				FE_LINK, uint(FE_GPU), dcgmErrorText(DCGM_ST_BADPARAM))
			require.EqualError(t, err, wantMessage)

			wrapped := fmt.Errorf("candidate link: %w", err)
			var nativeErr *Error
			require.ErrorAs(t, wrapped, &nativeErr)
			require.Equal(t, DCGM_ST_BADPARAM, int(nativeErr.Code))
			require.ErrorIs(t, wrapped, nativeErr)
		})
	}
}
