/*
 * Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeHeader(t *testing.T, contents string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "dcgm_fields.h")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing test header: %v", err)
	}
	return path
}

func findField(fields []Field, name string) (Field, bool) {
	for _, f := range fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}

// Plain integer defines are picked up and the aliases map is empty.
func TestParseHeader_IntegerDefines(t *testing.T) {
	path := writeHeader(t, `
/**
 * Field Foo.
 */
#define DCGM_FI_DEV_FOO 1

/**
 * Field Bar.
 */
#define DCGM_FI_DEV_BAR 2
`)

	fields, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	if len(fields) != 2 {
		t.Fatalf("want 2 fields, got %d: %+v", len(fields), fields)
	}
	if len(aliases) != 0 {
		t.Fatalf("want no aliases, got %v", aliases)
	}
}

// Alias inside #ifdef DCGM_DEPRECATED is recorded.
func TestParseHeader_AliasInsideDeprecatedBlock_Accepted(t *testing.T) {
	path := writeHeader(t, `
/**
 * Throughput.
 */
#define DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL 449

#define DCGM_DEPRECATED

#ifdef DCGM_DEPRECATED
#define DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL    DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL
#endif
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	target, ok := aliases["DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"]
	if !ok {
		t.Fatalf("expected BANDWIDTH_TOTAL alias to be recorded, got %v", aliases)
	}
	if target != "DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL" {
		t.Errorf("wrong alias target: %q", target)
	}
}

// Alias outside the #ifdef block is accepted when its preceding comment
// block contains "Deprecated:".
func TestParseHeader_AliasOutsideBlockWithDeprecatedComment_Accepted(t *testing.T) {
	path := writeHeader(t, `
/**
 * Clock events.
 */
#define DCGM_FI_DEV_CLOCKS_EVENT_REASONS 112

/**
 * Deprecated: Use DCGM_FI_DEV_CLOCKS_EVENT_REASONS instead
 */
#define DCGM_FI_DEV_CLOCK_THROTTLE_REASONS DCGM_FI_DEV_CLOCKS_EVENT_REASONS
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	target, ok := aliases["DCGM_FI_DEV_CLOCK_THROTTLE_REASONS"]
	if !ok {
		t.Fatalf("expected CLOCK_THROTTLE_REASONS alias to be recorded, got %v", aliases)
	}
	if target != "DCGM_FI_DEV_CLOCKS_EVENT_REASONS" {
		t.Errorf("wrong alias target: %q", target)
	}
}

// Alias outside the block with a non-deprecating comment (range sentinels,
// e.g. DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST) is silently rejected.
func TestParseHeader_AliasOutsideBlockWithoutDeprecation_Rejected(t *testing.T) {
	path := writeHeader(t, `
/**
 * Lane zero bytes.
 */
#define DCGM_FI_PROF_NVLINK_L0_TX_BYTES 1000

/**
 * NVLink throughput First.
 */
#define DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST DCGM_FI_PROF_NVLINK_L0_TX_BYTES
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	if _, ok := aliases["DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST"]; ok {
		t.Errorf("range sentinel was accepted as a deprecated alias: %v", aliases)
	}
	if len(aliases) != 0 {
		t.Errorf("expected no aliases, got %v", aliases)
	}
}

// The closing "*/" line must not be captured as comment content; the
// preceding comment should reach the #define intact. Pins the fix for the
// pre-existing "X represents /" rendering bug.
func TestParseHeader_ClosingStarSlashDoesNotPolluteComment(t *testing.T) {
	path := writeHeader(t, `
/**
 * Memory Application clocks
 */
#define DCGM_FI_DEV_APP_MEM_CLOCK 111
`)

	fields, _, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	f, ok := findField(fields, "DCGM_FI_DEV_APP_MEM_CLOCK")
	if !ok {
		t.Fatalf("field not parsed")
	}
	if f.Comment == "represents /" || strings.HasSuffix(f.Comment, "/") {
		t.Errorf("comment corrupted by */ line: %q", f.Comment)
	}
	if !strings.Contains(f.Comment, "Memory Application clocks") {
		t.Errorf("lost real comment content: %q", f.Comment)
	}
}

// A single-line "/** @} */" marker must not leave the parser stuck in
// comment-block mode. A following alias is handled normally.
func TestParseHeader_SingleLineBlockDoesNotGetStuck(t *testing.T) {
	path := writeHeader(t, `
/**
 * Lane zero.
 */
#define DCGM_FI_PROF_NVLINK_L0_TX_BYTES 1000

/** @} */

/**
 * NVLink throughput First.
 */
#define DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST DCGM_FI_PROF_NVLINK_L0_TX_BYTES
`)

	fields, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	// If the parser were stuck in comment-block mode, the L0_TX_BYTES field
	// would never be recorded because the "#define" handler wouldn't run.
	if _, ok := findField(fields, "DCGM_FI_PROF_NVLINK_L0_TX_BYTES"); !ok {
		t.Errorf("field after single-line block was lost: %+v", fields)
	}
	// And the non-deprecated alias after the marker should still be rejected.
	if _, ok := aliases["DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST"]; ok {
		t.Errorf("single-line marker caused false-positive on later alias: %v", aliases)
	}
}

// A "Deprecated:" comment attached to a numeric #define must not leak onto a
// later alias. The header shape under test mirrors DCGM_FI_DEV_PCIE_TX_THROUGHPUT.
func TestParseHeader_DeprecatedCommentOnNumericDefineDoesNotLeakToAlias(t *testing.T) {
	path := writeHeader(t, `
/**
 * Canonical target.
 */
#define DCGM_FI_PROF_PCIE_TX_BYTES 1010

/**
 * PCIe Tx utilization information
 *
 * Deprecated: Use DCGM_FI_PROF_PCIE_TX_BYTES instead.
 */
#define DCGM_FI_DEV_PCIE_TX_THROUGHPUT 200

/**
 * Lane zero.
 */
#define DCGM_FI_PROF_NVLINK_L0_TX_BYTES 1000

/**
 * NVLink throughput First.
 */
#define DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST DCGM_FI_PROF_NVLINK_L0_TX_BYTES
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	// The later alias must NOT be picked up as deprecated just because an
	// earlier numeric #define had a "Deprecated:" comment.
	if _, ok := aliases["DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST"]; ok {
		t.Errorf("deprecated state leaked from an earlier numeric define onto a later alias: %v", aliases)
	}
}

// resolveAliases returns an error when the target of a deprecated alias
// isn't in the fields slice.
func TestResolveAliases_TargetMissing(t *testing.T) {
	fields := []Field{{Name: "DCGM_FI_DEV_REAL", ID: 1}}
	aliases := map[string]string{"DCGM_FI_DEV_OLD": "DCGM_FI_DEV_GONE"}

	_, err := resolveAliases(fields, aliases)
	if err == nil {
		t.Fatalf("expected error when alias target is missing, got nil")
	}
	if !strings.Contains(err.Error(), "DCGM_FI_DEV_GONE") {
		t.Errorf("error should name the missing target, got %q", err)
	}
}

// Lowercase curated legacy entries in the existing output file are
// preserved across regeneration.
func TestGenerateOutput_PreservesCuratedLowercaseLegacy(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "const_fields.go")

	// Seed an existing output containing a curated lowercase entry.
	seeded := `package dcgm

const (
)
var dcgmFields = map[string]Short{
}
var legacyDCGMFields = map[string]Short{
	"dcgm_gpu_temp": 150,
}
`
	if err := os.WriteFile(outputPath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		t.Fatalf("extractLegacyFields: %v", err)
	}

	if got := legacyFields["dcgm_gpu_temp"]; got != 150 {
		t.Errorf("lost curated lowercase entry; got %v", legacyFields)
	}
}

// DCGM_FI_* uppercase entries in the existing output file are NOT
// preserved. They re-derive from the header via resolveAliases, so stale
// entries (aliases removed from the header) disappear naturally.
func TestGenerateOutput_StaleGeneratedAliasNotPreserved(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "const_fields.go")

	seeded := `package dcgm

const (
)
var dcgmFields = map[string]Short{
}
var legacyDCGMFields = map[string]Short{
	"dcgm_gpu_temp": 150,
	"DCGM_FI_DEV_SOMETHING_REMOVED": 99,
}
`
	if err := os.WriteFile(outputPath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		t.Fatalf("extractLegacyFields: %v", err)
	}

	if _, ok := legacyFields["DCGM_FI_DEV_SOMETHING_REMOVED"]; ok {
		t.Errorf("stale DCGM_FI_* entry was preserved; got %v", legacyFields)
	}
	if got := legacyFields["dcgm_gpu_temp"]; got != 150 {
		t.Errorf("curated lowercase entry lost alongside stale drop; got %v", legacyFields)
	}
}

// An uppercase non-DCGM_FI_* entry in the existing output is unrecognised
// provenance and causes a hard error.
func TestExtractLegacyFields_UnrecognisedUppercaseErrors(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "const_fields.go")

	seeded := `package dcgm
var legacyDCGMFields = map[string]Short{
	"MYSTERY_UPPERCASE_NAME": 777,
}
`
	if err := os.WriteFile(outputPath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	_, err := extractLegacyFields(outputPath)
	if err == nil {
		t.Fatalf("expected error on unrecognised uppercase entry, got nil")
	}
	if !strings.Contains(err.Error(), "MYSTERY_UPPERCASE_NAME") {
		t.Errorf("error should name the offending entry, got %q", err)
	}
}
