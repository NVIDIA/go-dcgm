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
// comment-block mode. Every assertion in this test describes something
// that happens strictly AFTER the marker -- if the parser got stuck,
// these assertions fail.
func TestParseHeader_SingleLineBlockDoesNotGetStuck(t *testing.T) {
	path := writeHeader(t, `
/** @} */

/**
 * Field described after the marker.
 */
#define DCGM_FI_DEV_AFTER_MARKER 999

#define DCGM_DEPRECATED

#ifdef DCGM_DEPRECATED
#define DCGM_FI_DEV_POST_MARKER_ALIAS DCGM_FI_DEV_AFTER_MARKER
#endif

/**
 * Range sentinel with a non-deprecating comment, after the marker.
 */
#define DCGM_FI_PROF_POST_MARKER_SENTINEL DCGM_FI_DEV_AFTER_MARKER
`)

	fields, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	// (1) If the parser were stuck in comment-block mode, the integer
	// define after the marker would never run the #define handler and
	// no field would be recorded.
	f, ok := findField(fields, "DCGM_FI_DEV_AFTER_MARKER")
	if !ok {
		t.Fatalf("parser stuck after /** @} */: field AFTER_MARKER was never recorded; got %+v", fields)
	}

	// (2) The marker line itself must not be captured as comment content
	// for the following field.
	if !strings.Contains(f.Comment, "Field described after the marker") {
		t.Errorf("comment on field AFTER_MARKER is not the expected one (marker line may have polluted state): %q", f.Comment)
	}

	// (3) The #ifdef DCGM_DEPRECATED block after the marker must still
	// be enterable -- the deprecated alias inside should be recorded.
	if _, ok := aliases["DCGM_FI_DEV_POST_MARKER_ALIAS"]; !ok {
		t.Errorf("deprecated alias defined after the marker was not recorded: %v", aliases)
	}

	// (4) And a non-deprecated alias after the marker should still be
	// silently rejected by the scope filter.
	if _, ok := aliases["DCGM_FI_PROF_POST_MARKER_SENTINEL"]; ok {
		t.Errorf("range sentinel after the marker was accepted as deprecated: %v", aliases)
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
func TestExtractLegacyFields_PreservesCuratedLowercaseEntries(t *testing.T) {
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
func TestExtractLegacyFields_StaleGeneratedAliasNotPreserved(t *testing.T) {
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

// A missing output file is normal on first-run; main() inspects the
// returned error with os.IsNotExist and starts with an empty legacy map.
func TestExtractLegacyFields_FileNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := extractLegacyFields(filepath.Join(dir, "does-not-exist.go"))
	if err == nil {
		t.Fatalf("expected error on missing output file, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("error should satisfy os.IsNotExist so main can tolerate it; got %v", err)
	}
}

// Mixed provenance in a single legacy map: lowercase preserved, DCGM_FI_*
// dropped (re-derived by resolveAliases elsewhere), no error.
func TestExtractLegacyFields_MixedEntriesReturnsOnlyLowercase(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "const_fields.go")

	seeded := `package dcgm
var legacyDCGMFields = map[string]Short{
	"dcgm_gpu_temp": 150,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL": 449,
	"dcgm_xid_errors": 230,
}
`
	if err := os.WriteFile(outputPath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		t.Fatalf("extractLegacyFields: %v", err)
	}
	if len(legacyFields) != 2 {
		t.Fatalf("want exactly 2 preserved entries, got %d: %v", len(legacyFields), legacyFields)
	}
	if legacyFields["dcgm_gpu_temp"] != 150 || legacyFields["dcgm_xid_errors"] != 230 {
		t.Errorf("lowercase entries not preserved correctly: %v", legacyFields)
	}
	if _, ok := legacyFields["DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"]; ok {
		t.Errorf("DCGM_FI_* entry should have been dropped: %v", legacyFields)
	}
}

// Bare alias inside #ifdef DCGM_DEPRECATED with no preceding comment: the
// block scope alone is sufficient. This pins the `inDeprecatedBlock`
// branch of the acceptance rule.
func TestParseHeader_AliasInsideDeprecatedBlockWithoutComment_Accepted(t *testing.T) {
	path := writeHeader(t, `
#define DCGM_FI_DEV_CANONICAL 42

#define DCGM_DEPRECATED

#ifdef DCGM_DEPRECATED
#define DCGM_FI_DEV_OLD_NAME DCGM_FI_DEV_CANONICAL
#endif
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}
	if target, ok := aliases["DCGM_FI_DEV_OLD_NAME"]; !ok {
		t.Errorf("bare alias inside deprecated block was dropped: %v", aliases)
	} else if target != "DCGM_FI_DEV_CANONICAL" {
		t.Errorf("wrong alias target: %q", target)
	}
}

// A nested #ifdef/#endif inside the deprecated block must not prematurely
// close it. An alias after the inner #endif but before the outer #endif
// is still in deprecated scope.
func TestParseHeader_NestedIfdefInsideDeprecatedBlock(t *testing.T) {
	path := writeHeader(t, `
#define DCGM_FI_DEV_CANONICAL 42

#define DCGM_DEPRECATED

#ifdef DCGM_DEPRECATED
#ifdef SOME_UNRELATED_GUARD
#define DCGM_FI_DEV_INNER_ALIAS DCGM_FI_DEV_CANONICAL
#endif
#define DCGM_FI_DEV_OUTER_ALIAS DCGM_FI_DEV_CANONICAL
#endif
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}
	// Both aliases should be recorded: the inner one while the outer guard
	// is active, the outer one after the inner #endif returns to the
	// deprecated block but not out of it.
	if _, ok := aliases["DCGM_FI_DEV_INNER_ALIAS"]; !ok {
		t.Errorf("alias inside nested #ifdef not recorded: %v", aliases)
	}
	if _, ok := aliases["DCGM_FI_DEV_OUTER_ALIAS"]; !ok {
		t.Errorf("alias after nested #endif but still inside deprecated block not recorded: %v", aliases)
	}
}

// The "deprecated:" marker match is case-insensitive. A header that
// capitalised the whole word should still trigger the heuristic.
func TestParseHeader_CaseInsensitiveDeprecatedMarker(t *testing.T) {
	path := writeHeader(t, `
#define DCGM_FI_DEV_CANONICAL 42

/**
 * DEPRECATED: Use DCGM_FI_DEV_CANONICAL instead.
 */
#define DCGM_FI_DEV_OLD_UPPER DCGM_FI_DEV_CANONICAL

/**
 * deprecated: lowercase form also counts.
 */
#define DCGM_FI_DEV_OLD_LOWER DCGM_FI_DEV_CANONICAL
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}
	if _, ok := aliases["DCGM_FI_DEV_OLD_UPPER"]; !ok {
		t.Errorf("uppercase DEPRECATED: did not trigger: %v", aliases)
	}
	if _, ok := aliases["DCGM_FI_DEV_OLD_LOWER"]; !ok {
		t.Errorf("lowercase deprecated: did not trigger: %v", aliases)
	}
}

// Adjectival uses of "deprecated" (without a trailing colon) must not
// trigger the heuristic. The pattern is specifically "deprecated:" so
// comments describing some other deprecation don't false-positive.
func TestParseHeader_DeprecatedWordWithoutColon_DoesNotTrigger(t *testing.T) {
	path := writeHeader(t, `
#define DCGM_FI_DEV_CANONICAL 42

/**
 * This replaces the now-deprecated DCGM_FI_DEV_FOO field; no colon here.
 */
#define DCGM_FI_DEV_ALIAS DCGM_FI_DEV_CANONICAL
`)

	_, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}
	if _, ok := aliases["DCGM_FI_DEV_ALIAS"]; ok {
		t.Errorf("adjectival 'deprecated' (no colon) false-positively triggered: %v", aliases)
	}
}

// Blank lines between the closing */ and its described #define must
// preserve the comment state so the comment still attaches to the field.
// This is the common real-world header layout.
func TestParseHeader_BlankLineBetweenCommentAndDefineAttaches(t *testing.T) {
	path := writeHeader(t, `
/**
 * Field with blank line between comment and define.
 */

#define DCGM_FI_DEV_CANONICAL 42

/**
 * Deprecated: with blank line before alias too.
 */

#define DCGM_FI_DEV_OLD DCGM_FI_DEV_CANONICAL
`)

	fields, aliases, err := parseHeader(path)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}

	f, ok := findField(fields, "DCGM_FI_DEV_CANONICAL")
	if !ok || !strings.Contains(f.Comment, "blank line between comment and define") {
		t.Errorf("comment did not attach across blank line: field=%+v", f)
	}
	if _, ok := aliases["DCGM_FI_DEV_OLD"]; !ok {
		t.Errorf("deprecated: marker did not carry across blank line: %v", aliases)
	}
}

// Full-pipeline integration: parseHeader -> resolveAliases ->
// extractLegacyFields -> generateOutput. Reads the emitted file and
// verifies the expected constants, canonical map entries, and legacy map
// entries all land in the right sections.
func TestGenerateOutput_FullPipeline(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "const_fields.go")

	// Seed a pre-existing output with a curated lowercase legacy entry
	// that the pipeline should preserve.
	seeded := `package dcgm
var legacyDCGMFields = map[string]Short{
	"dcgm_gpu_temp": 150,
}
`
	if err := os.WriteFile(outputPath, []byte(seeded), 0o644); err != nil {
		t.Fatalf("seed write: %v", err)
	}

	headerPath := writeHeader(t, `
/**
 * Canonical field.
 */
#define DCGM_FI_DEV_CANONICAL 42

#define DCGM_DEPRECATED

#ifdef DCGM_DEPRECATED
#define DCGM_FI_DEV_OLD_ALIAS DCGM_FI_DEV_CANONICAL
#endif
`)

	fields, aliases, err := parseHeader(headerPath)
	if err != nil {
		t.Fatalf("parseHeader: %v", err)
	}
	aliasLegacy, err := resolveAliases(fields, aliases)
	if err != nil {
		t.Fatalf("resolveAliases: %v", err)
	}
	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		t.Fatalf("extractLegacyFields: %v", err)
	}
	for name, id := range aliasLegacy {
		legacyFields[name] = id
	}

	if err := generateOutput(TemplateData{Fields: fields, LegacyFields: legacyFields}, outputPath); err != nil {
		t.Fatalf("generateOutput: %v", err)
	}

	out, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read emitted: %v", err)
	}
	got := string(out)

	// Typed constant.
	if !strings.Contains(got, "DCGM_FI_DEV_CANONICAL Short = 42") {
		t.Errorf("canonical typed constant missing in output:\n%s", got)
	}
	// dcgmFields canonical entry.
	if !strings.Contains(got, `"DCGM_FI_DEV_CANONICAL": 42`) {
		t.Errorf("canonical entry missing from dcgmFields:\n%s", got)
	}
	// legacyDCGMFields alias entry.
	if !strings.Contains(got, `"DCGM_FI_DEV_OLD_ALIAS": 42`) {
		t.Errorf("deprecated alias missing from legacyDCGMFields:\n%s", got)
	}
	// legacyDCGMFields curated lowercase entry.
	if !strings.Contains(got, `"dcgm_gpu_temp": 150`) {
		t.Errorf("curated lowercase entry not preserved:\n%s", got)
	}
}
