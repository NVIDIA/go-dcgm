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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIMEXHealthWatchConstant(t *testing.T) {
	const want HealthSystem = 0x2000

	require.Equal(t, want, DCGM_HEALTH_WATCH_IMEX)
	require.Equal(t, want, DCGM_HEALTH_WATCH_ALL&want)
}

func TestSystemWatchIMEX(t *testing.T) {
	require.Equal(t, "IMEX watches", systemWatch(DCGM_HEALTH_WATCH_IMEX))
}
