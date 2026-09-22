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

import "testing"

func TestSelectedFieldEntityGroupString(t *testing.T) {
	tests := []struct {
		group Field_Entity_Group
		want  string
	}{
		{FE_NONE, "unknown"},
		{FE_GPU, "GPU"},
		{FE_VGPU, "vGPU"},
		{FE_SWITCH, "NvSwitch"},
		{FE_GPU_I, "GPU Instance"},
		{FE_GPU_CI, "GPU Compute Instance"},
		{FE_LINK, "NvLink"},
		{FE_CPU, "CPU"},
		{FE_CPU_CORE, "CPU Core"},
		{FE_CONNECTX, "ConnectX"},
		{FE_COUNT, "unknown"},
		{Field_Entity_Group(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.group.String(); got != tt.want {
				t.Fatalf("Field_Entity_Group(%d).String() = %q, want %q", tt.group, got, tt.want)
			}
		})
	}
}
