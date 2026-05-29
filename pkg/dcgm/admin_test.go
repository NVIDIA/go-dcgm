/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandaloneConnectionArgsUsesV3ForSupportedConnectionStrings(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "vsock one arg", args: []string{"vsock://3:5555"}},
		{name: "vsock legacy socket flag ignored", args: []string{"vsock://3:5555", "0"}},
		{name: "uppercase vsock legacy socket flag ignored", args: []string{"VSOCK://3:5555", "0"}},
		{name: "tcp uri legacy socket flag ignored", args: []string{"tcp://127.0.0.1:5555", "0"}},
		{name: "uppercase tcp uri legacy socket flag ignored", args: []string{"TCP://127.0.0.1:5555", "0"}},
		{name: "unix uri legacy socket flag ignored", args: []string{"unix:///tmp/nv-hostengine", "1"}},
		{name: "uppercase unix uri legacy socket flag ignored", args: []string{"UNIX:///tmp/nv-hostengine", "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := standaloneConnectionArgs(tt.args...)
			require.NoError(t, err)
			require.True(t, got.useV3)
			require.Equal(t, tt.args[0], got.address)
			require.Empty(t, got.socketFlag)
		})
	}
}

func TestStandaloneConnectionArgsPreservesLegacyV2Arguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "host port", args: []string{"127.0.0.1:5555", "0"}},
		{name: "bracketed ipv6 host port", args: []string{"[::1]:5555", "0"}},
		{name: "unix socket legacy flag", args: []string{"/tmp/nv-hostengine", "1"}},
		{name: "unsupported http uri with legacy socket flag", args: []string{"http://127.0.0.1:5555", "0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := standaloneConnectionArgs(tt.args...)
			require.NoError(t, err)
			require.False(t, got.useV3)
			require.Equal(t, tt.args[0], got.address)
			require.Equal(t, tt.args[1], got.socketFlag)
		})
	}
}

func TestStandaloneConnectionArgsRejectsLegacyMissingSocketFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no args", args: nil},
		{name: "raw host without socket flag", args: []string{"127.0.0.1:5555"}},
		{name: "unsupported http uri without socket flag", args: []string{"http://127.0.0.1:5555"}},
		{name: "unsupported https uri without socket flag", args: []string{"https://127.0.0.1:5555"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := standaloneConnectionArgs(tt.args...)
			require.Error(t, err)
			require.Contains(t, err.Error(), "missing dcgm address and / or port")
		})
	}
}

func TestConnectStandaloneV3ReturnsErrorWhenSymbolUnavailable(t *testing.T) {
	oldLibHandle := dcgmLibHandle
	dcgmLibHandle = nil
	t.Cleanup(func() {
		dcgmLibHandle = oldLibHandle
	})

	err := connectStandaloneV3("vsock://3:5555")

	require.Error(t, err)
	require.Contains(t, err.Error(), "dcgmConnect_v3 is not available in libdcgm.so.4")
	require.Contains(t, err.Error(), "DCGM connection strings require DCGM 4.5.0 or newer")
}
