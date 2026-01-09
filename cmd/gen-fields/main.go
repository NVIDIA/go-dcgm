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
	fields, err := parseHeader(headerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing header: %v\n", err)
		os.Exit(1)
	}

	// Extract legacy fields from existing file
	legacyFields, err := extractLegacyFields(outputPath)
	if err != nil {
		// If file doesn't exist yet, start with empty legacy map
		legacyFields = make(map[string]int)
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

	fmt.Printf("Generated %d fields to %s\n", len(fields), outputPath)
}

func parseHeader(path string) ([]Field, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open header file: %w", err)
	}
	defer file.Close()

	// Pattern: #define DCGM_FI_XXX 123
	definePattern := regexp.MustCompile(`^#define\s+(DCGM_FI_\w+)\s+(\d+)`)
	commentPattern := regexp.MustCompile(`^\s*\*\s*(.+)$`)

	var fields []Field
	var lastComment string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for comments that describe the next field
		if strings.Contains(line, "/*") || strings.Contains(line, "*") {
			if matches := commentPattern.FindStringSubmatch(line); len(matches) > 1 {
				lastComment = strings.TrimSpace(matches[1])
			}
			continue
		}

		// Check for #define DCGM_FI_*
		if matches := definePattern.FindStringSubmatch(line); len(matches) == 3 {
			name := matches[1]
			idStr := matches[2]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			comment := lastComment
			if comment != "" {
				// Clean up comment
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
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading header file: %w", err)
	}

	// Sort by ID
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].ID < fields[j].ID
	})

	return fields, nil
}

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
				if err == nil {
					legacyFields[name] = id
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
