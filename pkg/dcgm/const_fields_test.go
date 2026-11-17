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

package dcgm

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"testing"
)

// TestFieldConstantsMatchFieldNameLookup checks that every DCGM_FI_* Go
// constant has the same ID as GetFieldID returns for the same name.
func TestFieldConstantsMatchFieldNameLookup(t *testing.T) {
	constants := fieldConstantValues(t)

	checked := 0
	for name, value := range constants {
		lookupID, ok := GetFieldID(name)
		if !ok {
			continue
		}
		checked++
		if int(lookupID) != value {
			t.Errorf("%s = %d, but GetFieldID(%q) = %d", name, value, name, lookupID)
		}
	}

	if checked == 0 {
		t.Fatal("no DCGM_FI_* constant has a field name lookup")
	}
}

// TestFieldConstantsMatchHeader checks that every numeric DCGM_FI_* define
// in dcgm_fields.h has a Go constant with the same ID.
func TestFieldConstantsMatchHeader(t *testing.T) {
	constants := fieldConstantValues(t)

	content, err := os.ReadFile("dcgm_fields.h")
	if err != nil {
		t.Fatalf("read dcgm_fields.h: %v", err)
	}

	definePattern := regexp.MustCompile(`(?m)^#define\s+(DCGM_FI_\w+)\s+(\d+)\b`)
	matches := definePattern.FindAllStringSubmatch(string(content), -1)
	if len(matches) == 0 {
		t.Fatal("no numeric DCGM_FI_* defines found in dcgm_fields.h")
	}

	for _, match := range matches {
		name := match[1]
		headerID, err := strconv.Atoi(match[2])
		if err != nil {
			t.Fatalf("parse %s value %q: %v", name, match[2], err)
		}

		value, ok := constants[name]
		if !ok {
			t.Errorf("%s is defined in dcgm_fields.h but has no Go constant", name)
			continue
		}
		if value != headerID {
			t.Errorf("%s = %d, but dcgm_fields.h defines it as %d", name, value, headerID)
		}
	}
}

// fieldConstantValues parses const_fields.go and returns the value of every
// DCGM_FI_* constant, following constants defined as another constant.
func fieldConstantValues(t *testing.T) map[string]int {
	t.Helper()

	file, err := parser.ParseFile(token.NewFileSet(), "const_fields.go", nil, 0)
	if err != nil {
		t.Fatalf("parse const_fields.go: %v", err)
	}

	// Collect each constant's right-hand side: a number or another constant.
	expressions := make(map[string]ast.Expr)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for i, name := range valueSpec.Names {
				expressions[name.Name] = valueSpec.Values[i]
			}
		}
	}

	values := make(map[string]int, len(expressions))
	for name := range expressions {
		values[name] = resolveFieldConstant(t, name, expressions, 0)
	}
	if len(values) == 0 {
		t.Fatal("no DCGM_FI_* constants found in const_fields.go")
	}

	return values
}

// resolveFieldConstant returns the numeric value of the named constant,
// following identifier references; depth stops a reference cycle.
func resolveFieldConstant(t *testing.T, name string, expressions map[string]ast.Expr, depth int) int {
	t.Helper()

	if depth > len(expressions) {
		t.Fatalf("constant %s refers to itself", name)
	}

	switch expr := expressions[name].(type) {
	case *ast.BasicLit:
		value, err := strconv.Atoi(expr.Value)
		if err != nil {
			t.Fatalf("parse %s value %q: %v", name, expr.Value, err)
		}
		return value
	case *ast.Ident:
		return resolveFieldConstant(t, expr.Name, expressions, depth+1)
	default:
		t.Fatalf("constant %s has an unsupported value %T", name, expr)
		return 0
	}
}
