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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type Field struct {
	Name    string
	ID      int
	Comment string
}

type TemplateData struct {
	Fields       []Field
	LegacyFields map[string]int
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: gen-fields <dcgm_fields.h> <const_fields.go>\n")
		os.Exit(1)
	}

	headerPath := os.Args[1]
	outputPath := os.Args[2]

	// Parse header file
	fields, aliases, err := parseHeader(headerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing header: %v\n", err)
		os.Exit(1)
	}

	// Resolve deprecated aliases to their target field IDs.
	aliasLegacy, err := resolveAliases(fields, aliases)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving aliases: %v\n", err)
		os.Exit(1)
	}

	// Extract legacy fields from existing file. Missing file on first-run is
	// fine; anything else (unreadable file, unrecognised legacy entry) is
	// treated as a hard error so we don't silently regenerate with lost
	// backward-compat names.
	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			legacyFields = make(map[string]int)
		} else {
			fmt.Fprintf(os.Stderr, "Error extracting legacy fields: %v\n", err)
			os.Exit(1)
		}
	}

	// Merge resolved aliases into the legacy map. Alias names start with
	// DCGM_FI_ and so never collide with the lowercase curated 1.x names.
	for name, id := range aliasLegacy {
		legacyFields[name] = id
	}

	// Generate output
	data := TemplateData{
		Fields:       fields,
		LegacyFields: legacyFields,
	}

	err = generateOutput(data, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %d fields (+ %d deprecated aliases) to %s\n",
		len(fields), len(aliasLegacy), outputPath)
}

// containsDeprecatedMarker reports whether the line contains the
// case-insensitive substring "deprecated:" -- the exact marker used in
// dcgm_fields.h to annotate deprecated aliases. Matching a loose "deprecated"
// substring would false-positive on adjectival mentions.
func containsDeprecatedMarker(line string) bool {
	return strings.Contains(strings.ToLower(line), "deprecated:")
}

func parseHeader(path string) ([]Field, map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open header file: %w", err)
	}
	defer file.Close()

	// #define DCGM_FI_XXX 123
	definePattern := regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+(\d+)`)
	// #define DCGM_FI_OLD DCGM_FI_NEW -- deprecated-alias shape.
	aliasPattern := regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+(DCGM_FI_\w+)\s*$`)
	// Content of a block-comment interior line: " * <content>".
	commentPattern := regexp.MustCompile(`^\s*\*\s*(.+)$`)

	var fields []Field
	aliases := make(map[string]string)

	var lastComment string
	// inCommentBlock tracks /** ... */ spans so the closing */ never feeds
	// commentPattern (which would otherwise capture "/" and corrupt
	// lastComment -- the origin of the "// X represents /" artefacts in the
	// previously-generated output).
	var inCommentBlock bool
	// commentHasDeprecated is set to true when any line inside the current
	// comment block contains case-insensitive "deprecated:". Consumed by the
	// alias handler to include outside-#ifdef-block cases like
	// DCGM_FI_DEV_CLOCK_THROTTLE_REASONS.
	var commentHasDeprecated bool
	// inDeprecatedBlock tracks the #ifdef DCGM_DEPRECATED ... #endif span.
	// A small nesting counter handles any nested #ifdef/#ifndef/#if inside.
	var inDeprecatedBlock bool
	var deprecatedBlockDepth int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Blank line preserves comment state.
		if trimmed == "" {
			continue
		}

		// Deprecated-block entry/exit. Track nesting so any #ifdef/#ifndef/#if
		// inside the block doesn't prematurely close it on its #endif.
		if trimmed == "#ifdef DCGM_DEPRECATED" {
			inDeprecatedBlock = true
			deprecatedBlockDepth = 1
			lastComment = ""
			commentHasDeprecated = false
			continue
		}
		if inDeprecatedBlock {
			if strings.HasPrefix(trimmed, "#ifdef ") ||
				strings.HasPrefix(trimmed, "#ifndef ") ||
				strings.HasPrefix(trimmed, "#if ") {
				deprecatedBlockDepth++
				lastComment = ""
				commentHasDeprecated = false
				continue
			}
			if strings.HasPrefix(trimmed, "#endif") {
				deprecatedBlockDepth--
				if deprecatedBlockDepth == 0 {
					inDeprecatedBlock = false
				}
				lastComment = ""
				commentHasDeprecated = false
				continue
			}
		}

		hasOpen := strings.Contains(line, "/*")
		hasClose := strings.Contains(line, "*/")

		// Single-line block like "/** @} */" or "/** Deprecated: X */":
		// inspect for the deprecated marker, do not enter block mode, do not
		// capture as field-describing content.
		if hasOpen && hasClose {
			if containsDeprecatedMarker(line) {
				commentHasDeprecated = true
			}
			inCommentBlock = false
			continue
		}

		// Block opener without a matching close on the same line.
		if hasOpen {
			inCommentBlock = true
			lastComment = ""
			commentHasDeprecated = false
			continue
		}

		// Block closer. Explicitly does NOT update lastComment; "*/" trimmed
		// would match commentPattern as " * /" and capture "/".
		if hasClose && inCommentBlock {
			inCommentBlock = false
			continue
		}

		// Interior of a block comment.
		if inCommentBlock {
			if matches := commentPattern.FindStringSubmatch(line); len(matches) > 1 {
				lastComment = strings.TrimSpace(matches[1])
			}
			if containsDeprecatedMarker(line) {
				commentHasDeprecated = true
			}
			continue
		}

		// Integer #define DCGM_FI_*.
		if matches := definePattern.FindStringSubmatch(line); len(matches) == 3 {
			name := matches[1]
			idStr := matches[2]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				lastComment = ""
				commentHasDeprecated = false
				continue
			}

			comment := lastComment
			if comment != "" {
				comment = strings.TrimSpace(comment)
				if !strings.HasPrefix(comment, "represents") {
					comment = "represents " + comment
				}
			}

			fields = append(fields, Field{
				Name:    name,
				ID:      id,
				Comment: comment,
			})

			lastComment = ""
			commentHasDeprecated = false
			continue
		}

		// Alias #define DCGM_FI_OLD DCGM_FI_NEW. Accept only if deprecated
		// either by position (inside #ifdef DCGM_DEPRECATED) or by an explicit
		// "Deprecated:" comment. Everything else is silently dropped -- the
		// common case is range sentinels like DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST
		// and it is not useful to log one per run.
		if matches := aliasPattern.FindStringSubmatch(line); len(matches) == 3 {
			aliasName := matches[1]
			targetName := matches[2]
			if inDeprecatedBlock || commentHasDeprecated {
				aliases[aliasName] = targetName
			}
			lastComment = ""
			commentHasDeprecated = false
			continue
		}

		// Any other non-blank line resets comment state so a comment meant for
		// one field never leaks onto a later construct.
		lastComment = ""
		commentHasDeprecated = false
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading header file: %w", err)
	}

	// Sort by ID for deterministic output.
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].ID < fields[j].ID
	})

	return fields, aliases, nil
}

// resolveAliases maps each deprecated alias to its target field's ID.
// Returns an error if any alias target is not a known field: we would rather
// fail generation loudly than silently ship a legacy map missing names that
// were previously exposed.
func resolveAliases(fields []Field, aliases map[string]string) (map[string]int, error) {
	fieldByName := make(map[string]int, len(fields))
	for _, f := range fields {
		fieldByName[f.Name] = f.ID
	}

	resolved := make(map[string]int, len(aliases))
	for alias, target := range aliases {
		id, ok := fieldByName[target]
		if !ok {
			return nil, fmt.Errorf(
				"deprecated alias %q points at unknown target %q; check dcgm_fields.h",
				alias, target)
		}
		resolved[alias] = id
	}
	return resolved, nil
}

// extractLegacyFields preserves curated legacy entries across regeneration.
// Provenance rules:
//
//   - Lowercase names (the hand-maintained DCGM 1.x backward-compat family,
//     e.g. "dcgm_gpu_temp") are preserved.
//   - DCGM_FI_* uppercase names are re-derived every run by resolveAliases
//     (see main), so they are skipped here to avoid stale entries persisting
//     after they disappear from the header.
//   - Any other uppercase name is an unrecognised provenance and fails
//     generation loudly, so hand-added names can't silently disappear on the
//     next regenerate and so unexpected patterns surface immediately.
func extractLegacyFields(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	legacyFields := make(map[string]int)

	// Pattern: "field_name": 123,
	entryPattern := regexp.MustCompile(`^\s*"([^"]+)":\s*(\d+),`)

	inLegacySection := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Look for the start of legacyDCGMFields map
		if strings.Contains(line, "var legacyDCGMFields") {
			inLegacySection = true
			continue
		}

		// If we're in the legacy section
		if inLegacySection {
			// Look for closing brace
			if strings.TrimSpace(line) == "}" {
				break
			}

			// Extract entries
			if matches := entryPattern.FindStringSubmatch(line); len(matches) == 3 {
				name := matches[1]
				id, err := strconv.Atoi(matches[2])
				if err != nil {
					continue
				}

				switch {
				case name == strings.ToLower(name):
					// Curated DCGM 1.x backward-compat name.
					legacyFields[name] = id
				case strings.HasPrefix(name, "DCGM_FI_"):
					// Generated alias; skip so it can be re-derived by
					// resolveAliases. Stale entries removed from the header
					// disappear naturally on regeneration.
					continue
				default:
					return nil, fmt.Errorf(
						"extractLegacyFields: %q has unrecognised provenance "+
							"(neither a lowercase DCGM 1.x name nor a DCGM_FI_* alias); "+
							"teach the generator how to regenerate it before adding it",
						name)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return legacyFields, nil
}

func generateOutput(data TemplateData, outputPath string) error {
	tmpl, err := template.New("fields").Parse(fileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
