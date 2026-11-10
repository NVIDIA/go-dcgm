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

package main

import "testing"

// TestCheckReleaseTag checks which tags each DCGM pin accepts.
func TestCheckReleaseTag(t *testing.T) {
	tests := []struct {
		name            string
		tag             string
		dcgmVersion     string
		packageRevision string
		wantErr         bool
	}{
		{name: "next patch", tag: "v1.4712.3", dcgmVersion: "4.7.1", packageRevision: "2"},
		{name: "first patch", tag: "v1.4712.0", dcgmVersion: "4.7.1", packageRevision: "2"},
		{name: "older pin", tag: "v1.4702.0", dcgmVersion: "4.7.0", packageRevision: "2"},
		{name: "revision 1", tag: "v1.4611.2", dcgmVersion: "4.6.1", packageRevision: "1"},
		{name: "two-digit minor", tag: "v1.41012.0", dcgmVersion: "4.10.1", packageRevision: "2"},
		{name: "other DCGM version", tag: "v1.4702.3", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "other package revision", tag: "v1.4711.3", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "no v prefix", tag: "1.4712.3", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "other major", tag: "v2.4712.3", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "no patch", tag: "v1.4712", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "suffix", tag: "v1.4712.3-rc1", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "dots are literal", tag: "v1x4712x3", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "bad DCGM version", tag: "v1.4712.3", dcgmVersion: "4.7.1-2", packageRevision: "2", wantErr: true},
		{name: "bad package revision", tag: "v1.4712.3", dcgmVersion: "4.7.1", packageRevision: "", wantErr: true},
		{name: "two-digit patch", tag: "v1.4712.10", dcgmVersion: "4.7.1", packageRevision: "2"},
		{name: "patch with leading zero", tag: "v1.4712.03", dcgmVersion: "4.7.1", packageRevision: "2", wantErr: true},
		{name: "DCGM version with leading zero", tag: "v1.40712.0", dcgmVersion: "4.07.1", packageRevision: "2", wantErr: true},
		{name: "package revision with leading zero", tag: "v1.47102.0", dcgmVersion: "4.7.1", packageRevision: "02", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkReleaseTag(tt.tag, tt.dcgmVersion, tt.packageRevision)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkReleaseTag(%q, %q, %q) error = %v, want error %t",
					tt.tag, tt.dcgmVersion, tt.packageRevision, err, tt.wantErr)
			}
		})
	}
}
