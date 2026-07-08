//go:build linux && cgo

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
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectedErrorMeta(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	meta := GetErrorMeta(DCGM_FR_PCI_REPLAY_RATE)
	require.NotNil(t, meta)
	require.Equal(t, DCGM_FR_PCI_REPLAY_RATE, meta.ErrorID)
	require.NotEmpty(t, meta.MessageFormat)
	require.NotEmpty(t, meta.Suggestion)
	// Keep this known entry exact so a runtime DCGM metadata change is detected.
	require.Equal(t, DCGM_ERROR_ISOLATE, meta.Severity)
	require.Equal(t, DCGM_FR_EC_HARDWARE_PCIE, meta.Category)

	require.Nil(t, GetErrorMeta(DCGM_FR_ERROR_SENTINEL))
	require.Nil(t, GetErrorMeta(HealthCheckErrorCode(math.MaxUint)))
}

func TestGoStringOrEmpty(t *testing.T) {
	require.Empty(t, goStringOrEmpty(nil))
}
