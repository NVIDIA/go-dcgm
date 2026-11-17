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
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const legacyFieldsCSVName = "legacy_fields.csv"

var dcgmConstNamePattern = regexp.MustCompile(`^DCGM_FI_[A-Z0-9_]+$`)

// Field describes a numeric field parsed from the DCGM header.
type Field struct {
	Name    string
	ID      int
	Comment string
}

// DeprecatedFieldAlias describes a deprecated DCGM field name that aliases a current field.
type DeprecatedFieldAlias struct {
	Name   string
	Target string
	ID     int
}

// TemplateData contains the complete input used to render const_fields.go.
type TemplateData struct {
	Fields            []Field
	DeprecatedAliases []DeprecatedFieldAlias
	LegacyFields      map[string]int
}

// main generates field constants or exits with a diagnostic.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run reads the header and compatibility data, then writes generated constants.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gen-fields", flag.ContinueOnError)
	flags.SetOutput(stderr)
	legacyFieldsFlag := flags.String(
		"legacy-fields",
		"",
		"CSV file containing curated legacy field names (default: output directory)",
	)
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if len(flags.Args()) != 2 {
		fmt.Fprintf(stderr, "Usage: gen-fields [--legacy-fields path] <dcgm_fields.h> <const_fields.go>\n")
		return 1
	}

	headerPath := flags.Arg(0)
	outputPath := flags.Arg(1)
	legacyFieldsPath := legacyFieldsCSVPath(*legacyFieldsFlag, outputPath)

	// Parse header file
	fields, aliases, err := parseHeader(headerPath)
	if err != nil {
		fmt.Fprintf(stderr, "Error parsing header: %v\n", err)
		return 1
	}

	// Resolve deprecated aliases to their target field IDs.
	deprecatedAliases, err := resolveDeprecatedFieldAliases(fields, aliases)
	if err != nil {
		fmt.Fprintf(stderr, "Error resolving aliases: %v\n", err)
		return 1
	}

	legacyFields, err := readLegacyFieldsCSV(legacyFieldsPath)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading legacy fields from %q: %v\n", legacyFieldsPath, err)
		return 1
	}
	curatedAliases, err := curatedFieldAliases(fields, deprecatedAliases, legacyFields)
	if err != nil {
		fmt.Fprintf(stderr, "Error checking legacy fields from %q: %v\n", legacyFieldsPath, err)
		return 1
	}

	// Merge resolved aliases into the legacy map. Curated names are lowercase
	// or DCGM field constant names.
	for _, alias := range deprecatedAliases {
		legacyFields[alias.Name] = alias.ID
	}
	deprecatedAliases = append(deprecatedAliases, curatedAliases...)
	sort.Slice(deprecatedAliases, func(i, j int) bool {
		if deprecatedAliases[i].ID != deprecatedAliases[j].ID {
			return deprecatedAliases[i].ID < deprecatedAliases[j].ID
		}
		return deprecatedAliases[i].Name < deprecatedAliases[j].Name
	})

	// Generate output
	data := TemplateData{
		Fields:            fields,
		DeprecatedAliases: deprecatedAliases,
		LegacyFields:      legacyFields,
	}

	err = generateOutput(data, outputPath)
	if err != nil {
		fmt.Fprintf(stderr, "Error generating output: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Generated %d fields (+ %d deprecated aliases) to %s\n",
		len(fields), len(deprecatedAliases), outputPath)
	return 0
}

// legacyFieldsCSVPath uses the explicit path or the compatibility list beside the output.
func legacyFieldsCSVPath(flagPath, outputPath string) string {
	if flagPath != "" {
		return flagPath
	}
	return filepath.Join(filepath.Dir(outputPath), legacyFieldsCSVName)
}

// containsDeprecatedMarker reports whether the line contains the
// case-insensitive substring "deprecated:" -- the exact marker used in
// dcgm_fields.h to annotate deprecated aliases. Matching a loose "deprecated"
// substring would false-positive on adjectival mentions.
func containsDeprecatedMarker(line string) bool {
	return strings.Contains(strings.ToLower(line), "deprecated:")
}

var (
	// Decimal or hexadecimal IDs, optionally wrapped in parentheses.
	numericFieldPattern = regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+([()\s]*(?:0[xX][0-9a-fA-F]+|[0-9]+)[()\s]*)(?:\s*(?://.*|/\*.*\*/))?\s*$`)
	// Renamed aliases may have a trailing comment containing their numeric ID.
	fieldAliasPattern         = regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+(DCGM_FI_\w+)(?:\s*(?://.*|/\*.*\*/))?\s*$`)
	fieldExpressionPattern    = regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+[\w\s()+*/-]+$`)
	inlineBlockCommentPattern = regexp.MustCompile(`/\*.*?\*/`)
	fieldCommentPattern       = regexp.MustCompile(`^\s*\*\s*(.+)$`)
)

type headerParser struct {
	fields               []Field
	aliases              map[string]string
	skippedExpressions   map[string]bool
	lastComment          string
	inCommentBlock       bool
	commentHasDeprecated bool
	inDeprecatedBlock    bool
	deprecatedBlockDepth int
}

// parseHeader returns numeric fields sorted by ID and deprecated symbolic aliases.
// Arithmetic definitions and aliases referring to them are excluded after scanning.
func parseHeader(path string) ([]Field, map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open header file: %w", err)
	}
	defer file.Close()

	parser := headerParser{aliases: make(map[string]string), skippedExpressions: make(map[string]bool)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := parser.parseLine(scanner.Text()); err != nil {
			return nil, nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading header file: %w", err)
	}

	// Check after scanning so aliases may precede their expression targets.
	for alias, target := range parser.aliases {
		if parser.skippedExpressions[target] {
			delete(parser.aliases, alias)
		}
	}
	sort.Slice(parser.fields, func(i, j int) bool { return parser.fields[i].ID < parser.fields[j].ID })
	return parser.fields, parser.aliases, nil
}

// parseLine processes one header line, preserving comment and guard state between calls.
func (p *headerParser) parseLine(line string) error {
	trimmed := strings.TrimSpace(line)
	// Blank lines preserve comment state.
	if trimmed == "" {
		return nil
	}
	if p.parseDeprecatedDirective(trimmed) {
		return nil
	}
	if p.parseCommentLine(line, trimmed) {
		return nil
	}
	// Remove complete block comments before classifying definitions; code after them still counts.
	line = inlineBlockCommentPattern.ReplaceAllString(line, " ")
	if matched, err := p.parseNumericDefinition(line); matched {
		return err
	}
	if p.parseAliasDefinition(line) {
		return nil
	}
	p.recordSkippedExpression(line)
	// Unrelated constructs must not inherit a preceding field's comment.
	p.clearComment()
	return nil
}

// clearComment prevents a field's documentation from leaking into a later definition.
func (p *headerParser) clearComment() {
	p.lastComment = ""
	p.commentHasDeprecated = false
}

// parseDeprecatedDirective tracks nested guards within a deprecated block.
// Its result tells the caller whether the directive was handled.
func (p *headerParser) parseDeprecatedDirective(line string) bool {
	if !p.inDeprecatedBlock && isDeprecatedBlockStart(line) {
		p.inDeprecatedBlock = true
		p.deprecatedBlockDepth = 1
		p.clearComment()
		return true
	}
	if !p.inDeprecatedBlock {
		return false
	}
	// Nested guards must not close the surrounding deprecated block early.
	if strings.HasPrefix(line, "#ifdef ") || strings.HasPrefix(line, "#ifndef ") || strings.HasPrefix(line, "#if ") {
		p.deprecatedBlockDepth++
		p.clearComment()
		return true
	}
	if strings.HasPrefix(line, "#endif") {
		p.deprecatedBlockDepth--
		if p.deprecatedBlockDepth == 0 {
			p.inDeprecatedBlock = false
		}
		p.clearComment()
		return true
	}
	return false
}

// noteDeprecated retains a deprecation marker across lines of the current comment.
func (p *headerParser) noteDeprecated(line string) {
	if containsDeprecatedMarker(line) {
		p.commentHasDeprecated = true
	}
}

// parseCommentLine collects comment lines, leaving trailing comments for definition parsing.
func (p *headerParser) parseCommentLine(line, trimmed string) bool {
	if strings.HasPrefix(trimmed, "//") {
		p.noteDeprecated(line)
		return true
	}
	hasOpen := strings.Contains(line, "/*")
	hasClose := strings.Contains(line, "*/")
	if hasOpen && hasClose {
		if strings.HasPrefix(trimmed, "/*") {
			p.lastComment = ""
			p.commentHasDeprecated = containsDeprecatedMarker(line)
			p.inCommentBlock = false
			return true
		}
		// Inline block comments stay on the line for definition parsing.
		p.noteDeprecated(line)
	}
	if hasOpen && !hasClose {
		p.inCommentBlock = true
		p.clearComment()
		return true
	}
	if hasClose && p.inCommentBlock {
		// Never capture the closing slash as a field description.
		p.inCommentBlock = false
		return true
	}
	if p.inCommentBlock {
		if content := fieldCommentText(line); content != "" {
			p.lastComment = content
		}
		p.noteDeprecated(line)
		return true
	}
	return false
}

// fieldCommentText extracts the description from the comment preceding a field definition.
func fieldCommentText(line string) string {
	matches := fieldCommentPattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}
	content := strings.TrimSpace(matches[1])
	if strings.HasPrefix(content, "@") || strings.HasPrefix(content, `\`) {
		return ""
	}
	return content
}

// parseNumericFieldID parses a field identifier without accepting expressions.
func parseNumericFieldID(value string) (int, error) {
	value = strings.TrimSpace(value)
	// Peel complete pairs only; malformed parentheses still fail numeric parsing.
	for strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
		value = strings.TrimSpace(value[1 : len(value)-1])
	}
	// Leading zeroes on decimal IDs do not select octal notation.
	base := 10
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		base = 16
		value = value[2:]
	}
	id, err := strconv.ParseInt(value, base, strconv.IntSize)
	return int(id), err
}

// parseNumericDefinition records a literal field ID with its pending documentation.
// Invalid numeric IDs abort parsing rather than silently dropping a field.
func (p *headerParser) parseNumericDefinition(line string) (bool, error) {
	matches := numericFieldPattern.FindStringSubmatch(line)
	if len(matches) != 3 {
		return false, nil
	}
	id, err := parseNumericFieldID(matches[2])
	if err != nil {
		return true, fmt.Errorf("invalid numeric field %s with ID %q: %w", matches[1], matches[2], err)
	}
	comment := strings.TrimSpace(p.lastComment)
	if comment != "" && !strings.HasPrefix(comment, "represents") {
		comment = "represents " + comment
	}
	p.fields = append(p.fields, Field{Name: matches[1], ID: id, Comment: generatedFieldComment(matches[1], comment)})
	p.clearComment()
	return true, nil
}

// parseAliasDefinition records deprecated symbolic aliases and omits ordinary aliases.
func (p *headerParser) parseAliasDefinition(line string) bool {
	matches := fieldAliasPattern.FindStringSubmatch(line)
	if len(matches) != 3 {
		return false
	}
	if p.inDeprecatedBlock || p.commentHasDeprecated {
		p.aliases[matches[1]] = matches[2]
	}
	p.clearComment()
	return true
}

// recordSkippedExpression marks arithmetic macros so their aliases are excluded too.
func (p *headerParser) recordSkippedExpression(line string) {
	// Comment punctuation must not turn an unrelated macro into arithmetic.
	definition, _, _ := strings.Cut(line, "//")
	if matches := fieldExpressionPattern.FindStringSubmatch(definition); len(matches) == 2 && strings.ContainsAny(definition, "+*/-") {
		p.skippedExpressions[matches[1]] = true
	}
}

// generatedFieldComment replaces known incomplete upstream comments and keeps others.
func generatedFieldComment(name, comment string) string {
	switch name {
	case "DCGM_FI_PROF_NVLINK_TX_BYTES_PER_SECOND":
		return "represents aggregate NVLink transmit bytes per second for a GPU. For a DCGM_FE_LINK entity, dcgm_link_t selects the GPU and link whose traffic is reported."
	case "DCGM_FI_PROF_NVLINK_RX_BYTES_PER_SECOND":
		return "represents aggregate NVLink receive bytes per second for a GPU. For a DCGM_FE_LINK entity, dcgm_link_t selects the GPU and link whose traffic is reported."
	default:
		return comment
	}
}

// isDeprecatedBlockStart recognizes supported spellings of DCGM's deprecated guard.
func isDeprecatedBlockStart(trimmed string) bool {
	switch trimmed {
	case "#ifdef DCGM_DEPRECATED",
		"#if DCGM_DEPRECATED",
		"#if defined(DCGM_DEPRECATED)",
		"#if defined DCGM_DEPRECATED":
		return true
	default:
		return false
	}
}

// resolveDeprecatedFieldAliases maps each deprecated alias to its target field's ID.
// Returns an error if any alias target is not a known field: we would rather
// fail generation loudly than silently ship a legacy map missing names that
// were previously exposed.
func resolveDeprecatedFieldAliases(fields []Field, aliases map[string]string) ([]DeprecatedFieldAlias, error) {
	fieldByName := make(map[string]int, len(fields))
	for _, f := range fields {
		fieldByName[f.Name] = f.ID
	}

	resolved := make([]DeprecatedFieldAlias, 0, len(aliases))
	for alias, target := range aliases {
		id, ok := fieldByName[target]
		if !ok {
			return nil, fmt.Errorf(
				"deprecated alias %q points at unknown target %q; check dcgm_fields.h",
				alias, target,
			)
		}
		resolved = append(resolved, DeprecatedFieldAlias{
			Name:   alias,
			Target: target,
			ID:     id,
		})
	}

	sort.Slice(resolved, func(i, j int) bool {
		if resolved[i].ID != resolved[j].ID {
			return resolved[i].ID < resolved[j].ID
		}
		return resolved[i].Name < resolved[j].Name
	})

	return resolved, nil
}

// readLegacyFieldsCSV loads field definitions retained for compatibility.
func readLegacyFieldsCSV(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open legacy fields CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comment = '#'
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read legacy fields CSV header: %w", err)
	}
	if len(header) != 2 || strings.TrimSpace(header[0]) != "name" || strings.TrimSpace(header[1]) != "id" {
		return nil, errors.New("legacy fields CSV header must be: name,id")
	}

	legacyFields := make(map[string]int)
	for row := 2; ; row++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("legacy fields CSV row %d: %w", row, err)
		}
		if len(record) != 2 {
			return nil, fmt.Errorf("legacy fields CSV row %d: expected 2 columns, got %d", row, len(record))
		}

		name := strings.TrimSpace(record[0])
		idText := strings.TrimSpace(record[1])
		if name == "" {
			return nil, fmt.Errorf("legacy fields CSV row %d: name is required", row)
		}
		if name != strings.ToLower(name) && !dcgmConstNamePattern.MatchString(name) {
			return nil, fmt.Errorf("legacy fields CSV row %d: %q must be lowercase or a DCGM field constant", row, name)
		}
		if idText == "" {
			return nil, fmt.Errorf("legacy fields CSV row %d: id is required", row)
		}
		id, err := strconv.Atoi(idText)
		if err != nil || id < 0 {
			return nil, fmt.Errorf("legacy fields CSV row %d: %q is not a valid non-negative integer ID", row, idText)
		}
		if _, exists := legacyFields[name]; exists {
			return nil, fmt.Errorf("legacy fields CSV row %d: duplicate name %q", row, name)
		}
		legacyFields[name] = id
	}

	return legacyFields, nil
}

// curatedFieldAliases returns a Go constant alias for each DCGM_FI_* name in
// the curated legacy CSV that dcgm_fields.h no longer defines, pointing at the
// header field with the curated ID. A curated name that the header still
// defines, as a field or a deprecated alias, must carry the header's ID: the
// header decides what a DCGM field name means, so a mismatch fails generation
// instead of giving a Go constant a different ID than GetFieldID returns for
// the same name.
func curatedFieldAliases(
	fields []Field,
	headerAliases []DeprecatedFieldAlias,
	legacyFields map[string]int,
) ([]DeprecatedFieldAlias, error) {
	headerIDs := make(map[string]int, len(fields)+len(headerAliases))
	fieldNamesByID := make(map[int]string, len(fields))
	for _, field := range fields {
		headerIDs[field.Name] = field.ID
		fieldNamesByID[field.ID] = field.Name
	}
	for _, alias := range headerAliases {
		headerIDs[alias.Name] = alias.ID
	}

	var aliases []DeprecatedFieldAlias
	for name, id := range legacyFields {
		if !dcgmConstNamePattern.MatchString(name) {
			continue
		}

		// A name the header defines keeps the header's ID.
		if headerID, defined := headerIDs[name]; defined {
			if headerID != id {
				return nil, fmt.Errorf("%q has ID %d, but dcgm_fields.h defines it as %d", name, id, headerID)
			}
			continue
		}

		target, ok := fieldNamesByID[id]
		if !ok {
			return nil, fmt.Errorf("%q has ID %d, which dcgm_fields.h does not define", name, id)
		}
		aliases = append(aliases, DeprecatedFieldAlias{Name: name, Target: target, ID: id})
	}

	sort.Slice(aliases, func(i, j int) bool { return aliases[i].Name < aliases[j].Name })

	return aliases, nil
}

// generateOutput renders the field constants and their descriptions as Go source.
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
