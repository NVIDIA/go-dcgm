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

const fileTemplate = `/*
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

package dcgm

const (
{{- range .Fields}}
	// {{.Name}}{{if .Comment}} {{.Comment}}{{end}}
	{{.Name}} Short = {{.ID}}
{{- end}}
)

// dcgmFields maps field names to their IDs
var dcgmFields = map[string]Short{
{{- range .Fields}}
	"{{.Name}}": {{.ID}},
{{- end}}
}

// legacyDCGMFields maps legacy field names to their IDs
var legacyDCGMFields = map[string]Short{
{{- range $name, $id := .LegacyFields}}
	"{{$name}}": {{$id}},
{{- end}}
}

// GetFieldID returns the DCGM field ID for a given field name and whether it was found
// It first checks the current field IDs, then falls back to legacy field IDs if not found
func GetFieldID(fieldName string) (Short, bool) {
	// First check current field IDs
	if fieldID, ok := dcgmFields[fieldName]; ok {
		return fieldID, true
	}

	// Then check legacy field IDs
	if fieldID, ok := legacyDCGMFields[fieldName]; ok {
		return fieldID, true
	}

	return 0, false
}

// GetFieldIDOrPanic returns the DCGM field ID for a given field name
// It panics if the field name is not found in either current or legacy maps
func GetFieldIDOrPanic(fieldName string) Short {
	fieldID, ok := GetFieldID(fieldName)
	if !ok {
		panic("field name not found: " + fieldName)
	}
	return fieldID
}

// IsLegacyField returns true if the given field name is a legacy field
func IsLegacyField(fieldName string) bool {
	_, ok := legacyDCGMFields[fieldName]
	return ok
}

// IsCurrentField returns true if the given field name is a current field
func IsCurrentField(fieldName string) bool {
	_, ok := dcgmFields[fieldName]
	return ok
}
`
