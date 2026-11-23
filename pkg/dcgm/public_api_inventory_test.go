package dcgm

import (
	"encoding/csv"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
)

var (
	publicDeclaration = regexp.MustCompile(`(?s)[^;{}]*\bDCGM_PUBLIC_API\b[^;{}]*;`)
	cFunctionName     = regexp.MustCompile(`\b([A-Za-z_]\w*)\s*\(`)
	cLiterals         = regexp.MustCompile(`(?s)"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'`)
	cTokens           = regexp.MustCompile(cLiterals.String() + `|/\*.*?\*/|//[^\n]*`)
)

// publicAPINames lists exported C declarations while ignoring comments.
func publicAPINames(header string) []string {
	names := []string{}
	for _, declaration := range publicDeclaration.FindAllString(stripCComments(header), -1) {
		match := cFunctionName.FindStringSubmatch(declaration)
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}
	sort.Strings(names)
	return names
}

// inventoryProblems checks every public declaration and inventory row for matching evidence.
func inventoryProblems(names []string, entries map[string][]string, calls map[string]bool) []string {
	problems := []string{}
	declared := map[string]bool{}
	for _, name := range names {
		declared[name] = true
		entry, ok := entries[name]
		if !ok {
			problems = append(problems, "unclassified public API: "+name)
			continue
		}
		switch entry[0] {
		case "bound":
			if !calls[name] {
				problems = append(problems, "bound API has no native call evidence: "+name)
			}
		case "not bound yet", "excluded":
			if strings.TrimSpace(entry[1]) == "" {
				problems = append(problems, "API needs a reason: "+name)
			}
		default:
			problems = append(problems, "invalid API classification: "+name)
		}
	}
	for name := range entries {
		if !declared[name] {
			problems = append(problems, "inventory API absent from header: "+name)
		}
	}
	sort.Strings(problems)
	return problems
}

// dynamicAPIBridge describes optional APIs resolved and called through a typed pointer.
type dynamicAPIBridge struct {
	symbol, goFile, cFile, bridge, functionType, parameters, arguments string
}

var dynamicAPIBridges = []dynamicAPIBridge{
	{
		"dcgmVersionInfo", "version_info.go", "version_info.go", "loadedClientVersion", "clientVersionFn",
		`dcgmVersionInfo_t\s*\*`, "info",
	},
	{
		"dcgmHostengineVersionInfo", "version_info.go", "version_info.go", "loadedHostengineVersion", "hostengineVersionFn",
		`dcgmHandle_t\s*,\s*dcgmVersionInfo_t\s*\*`, "handle,info",
	},
	{
		"dcgmPolicyImport_v2", "field_policy.go", "field_policy_import.h", "callPolicyImportV2", "policyImportV2Function",
		`dcgmHandle_t\s*,\s*const\s+char\s*\*\s*,\s*unsigned\s+int\s*\*\s*,\s*unsigned\s+int\s*\*`, "handle,path,parsed,skipped",
	},
	{
		"dcgmConfigReset", "config.go", "config.go", "go_dcgmConfigReset", "goConfigResetFn",
		`dcgmHandle_t\s*,\s*dcgmGpuGrp_t\s*,\s*dcgmStatus_t`, "handle,group,status",
	},
}

// inventoryNativeCalls collects real Go calls and verifies the required typed native bridges.
func inventoryNativeCalls(sources map[string]string) (map[string]bool, error) {
	calls := map[string]bool{}
	for path, source := range sources {
		if filepath.Ext(path) != ".go" {
			continue
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
		if err != nil {
			return nil, fmt.Errorf("parse %s for native calls: %w", path, err)
		}
		for name := range goNativeCalls(file) {
			calls[name] = true
		}
	}
	for _, bridge := range dynamicAPIBridges {
		// Optional APIs require their actual bridge; a stray direct call is insufficient.
		bound, err := dynamicBridgePresent(sources, bridge)
		if err != nil {
			return nil, err
		}
		calls[bridge.symbol] = bound
	}

	bound, err := pidInfoBridgePresent(sources["process_info.go"])
	if err != nil {
		return nil, err
	}
	calls["dcgmGetPidInfo"] = bound
	return calls, nil
}

// pidInfoBridgePresent requires a Go call to the typed PID bridge and its exact native conversion.
func pidInfoBridgePresent(source string) (bool, error) {
	if source == "" {
		return false, nil
	}

	file, err := parser.ParseFile(token.NewFileSet(), "process_info.go", source, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("parse process_info.go for PID bridge: %w", err)
	}
	if !goCallsCBridge(file, "getPidInfoV2") {
		return false, nil
	}

	preamble := stripCComments(cgoPreamble(file))
	preamble = maskCLiterals(preamble)
	bridge := regexp.MustCompile(`\bstatic\s+inline\s+dcgmReturn_t\s+getPidInfoV2\s*\(\s*dcgmHandle_t\s+dcgmHandle\s*,\s*dcgmGpuGrp_t\s+groupId\s*,\s*dcgmPidInfo_v2\s*\*\s*pidInfo\s*\)\s*\{\s*return\s+dcgmGetPidInfo\s*\(\s*dcgmHandle\s*,\s*groupId\s*,\s*\(\s*dcgmPidInfo_t\s*\*\s*\)\s*pidInfo\s*\)\s*;\s*\}`)
	return bridge.MatchString(preamble), nil
}

// dynamicBridgePresent requires the Go call, included C bridge and typed symbol invocation.
func dynamicBridgePresent(sources map[string]string, bridge dynamicAPIBridge) (bool, error) {
	source, ok := sources[bridge.goFile]
	if !ok {
		return false, nil
	}

	file, err := parser.ParseFile(token.NewFileSet(), bridge.goFile, source, parser.ParseComments)
	if err != nil {
		return false, fmt.Errorf("parse %s for native bridge: %w", bridge.goFile, err)
	}
	if !goCallsCBridge(file, bridge.bridge) {
		return false, nil
	}

	preamble := cgoPreamble(file)
	cSource := preamble
	if bridge.cFile != bridge.goFile {
		include := regexp.MustCompile(`#include\s*"` + regexp.QuoteMeta(bridge.cFile) + `"`)
		if !include.MatchString(stripCComments(preamble)) {
			return false, nil
		}
		cSource = sources[bridge.cFile]
	}

	return typedDynamicBridge(cSource, bridge), nil
}

// goCallsCBridge reports whether the Go syntax tree calls the named C bridge.
func goCallsCBridge(file *ast.File, bridge string) bool {
	return goNativeCalls(file)[bridge]
}

// goNativeCalls finds C calls in Go syntax, excluding comment and string decoys.
func goNativeCalls(file *ast.File) map[string]bool {
	calls := map[string]bool{}
	if !slices.ContainsFunc(file.Imports, func(spec *ast.ImportSpec) bool { return spec.Path.Value == `"C"` }) {
		return calls
	}
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		receiver, ok := selector.X.(*ast.Ident)
		if ok && receiver.Name == "C" && receiver.Obj == nil {
			calls[selector.Sel.Name] = true
		}

		return true
	})
	return calls
}

// cgoPreamble returns the comment attached to import C, where native bridges are declared.
func cgoPreamble(file *ast.File) string {
	for _, declaration := range file.Decls {
		imports, ok := declaration.(*ast.GenDecl)
		if !ok || imports.Tok != token.IMPORT || imports.Doc == nil {
			continue
		}
		for _, spec := range imports.Specs {
			if spec.(*ast.ImportSpec).Path.Value == `"C"` {
				return imports.Doc.Text()
			}
		}
	}

	return ""
}

// stripCComments masks C comments while preserving literals and source offsets.
func stripCComments(source string) string {
	return cTokens.ReplaceAllStringFunc(source, func(token string) string {
		if strings.HasPrefix(token, "/") {
			return strings.Repeat(" ", len(token))
		}

		return token
	})
}

// maskCLiterals masks C literals without changing source offsets.
func maskCLiterals(source string) string {
	return cLiterals.ReplaceAllStringFunc(source, func(literal string) string {
		return strings.Repeat(" ", len(literal))
	})
}

// typedDynamicBridge requires a typed resolver and invocation inside the same native bridge.
func typedDynamicBridge(source string, bridge dynamicAPIBridge) bool {
	source = stripCComments(source)
	masked := maskCLiterals(source)
	typeName := regexp.QuoteMeta(bridge.functionType)
	typedef := regexp.MustCompile(`\btypedef\s+dcgmReturn_t\s*\(\s*\*\s*` + typeName + `\s*\)\s*\(\s*` + bridge.parameters + `\s*\)\s*;`)
	if !typedef.MatchString(masked) {
		return false
	}

	definition := regexp.MustCompile(`\b` + regexp.QuoteMeta(bridge.bridge) + `\s*\([^;{}]*\)\s*\{`)
	start := definition.FindStringIndex(masked)
	if start == nil {
		return false
	}

	depth, end := 1, start[1]
	for end < len(masked) && depth > 0 {
		switch masked[end] {
		case '{':
			depth++
		case '}':
			depth--
		}
		end++
	}
	if depth != 0 {
		return false
	}

	body := source[start[1] : end-1]
	resolver := regexp.MustCompile(`\b` + typeName + `\s+(\w+)\s*=\s*\(\s*` + typeName + `\s*\)\s*dlsym\s*\(\s*library\s*,\s*"` + regexp.QuoteMeta(bridge.symbol) + `"\s*\)\s*;`)
	maskedBody := masked[start[1] : end-1]
	args := strings.ReplaceAll(regexp.QuoteMeta(bridge.arguments), ",", `\s*,\s*`)
	for _, resolved := range resolver.FindAllStringSubmatchIndex(body, -1) {
		// Keep the dlsym symbol string, but require its declaration to be code.
		if maskedBody[resolved[0]:resolved[2]] != body[resolved[0]:resolved[2]] {
			continue
		}

		name := body[resolved[2]:resolved[3]]
		invocation := regexp.MustCompile(`\breturn\s+` + regexp.QuoteMeta(name) + `\s*\(\s*` + args + `\s*\)\s*;`)
		if invocation.MatchString(maskedBody) {
			return true
		}
	}

	return false
}

// readPublicAPIEntries reads unique API classifications and reasons from the inventory CSV.
func readPublicAPIEntries(path string) (map[string][]string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(strings.NewReader(string(contents))).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty API inventory: %s", path)
	}

	entries := map[string][]string{}
	for _, record := range records[1:] {
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid inventory entry: %q", record)
		}
		if _, ok := entries[record[0]]; ok {
			return nil, fmt.Errorf("duplicate inventory API: %s", record[0])
		}
		entries[record[0]] = record[1:]
	}

	return entries, nil
}

// readInventorySources loads binding sources while excluding tests and test fixtures.
func readInventorySources() (map[string]string, error) {
	sources := map[string]string{}
	for _, pattern := range []string{"*.go", "*.c", "field_policy_import.h"} {
		paths, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		for _, path := range paths {
			if strings.HasSuffix(path, "_test.go") || strings.Contains(path, "test_helpers") {
				continue
			}

			source, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil, readErr
			}
			sources[path] = string(source)
		}
	}

	return sources, nil
}

// TestPublicAPIInventory covers dcgm_agent.h. The unsupported test/internal
// header dcgm_test_apis.h is excluded; see testdata/README.md for the scope.
func TestPublicAPIInventory(t *testing.T) {
	header, err := os.ReadFile("dcgm_agent.h")
	if err != nil {
		t.Fatal(err)
	}

	entries, err := readPublicAPIEntries("testdata/public_api.csv")
	if err != nil {
		t.Fatal(err)
	}

	sources, err := readInventorySources()
	if err != nil {
		t.Fatal(err)
	}

	calls, err := inventoryNativeCalls(sources)
	if err != nil {
		t.Fatal(err)
	}
	if problems := inventoryProblems(publicAPINames(string(header)), entries, calls); len(problems) > 0 {
		t.Fatal(strings.Join(problems, "\n"))
	}
}

// TestPublicAPIDeclarationShapes checks declaration parsing and inventory coverage in both directions.
func TestPublicAPIDeclarationShapes(t *testing.T) {
	header := `/* DCGM_PUBLIC_API int fake(); */
 dcgmReturn_t DCGM_PUBLIC_API dcgmOne(void);
 dcgmReturn_t
 DCGM_PUBLIC_API dcgmTwo(void);
 DCGM_PUBLIC_API const char *errorString(int result);
 DCGM_PUBLIC_API dcgmReturn_t otherName(int id);`
	got := fmt.Sprint(publicAPINames(header))
	if got != "[dcgmOne dcgmTwo errorString otherName]" {
		t.Fatal(got)
	}

	entries := map[string][]string{"dcgmOne": {"bound", ""}, "dcgmTwo": {"excluded", "replaced"}, "errorString": {"bound", ""}, "otherName": {"not bound yet", "future work"}}
	calls := map[string]bool{"dcgmOne": true, "errorString": true}
	added := publicAPINames(header + "\ndcgmReturn_t DCGM_PUBLIC_API newFunction(void);")
	if got := inventoryProblems(added, entries, calls); len(got) != 1 || got[0] != "unclassified public API: newFunction" {
		t.Fatal(got)
	}
	delete(calls, "dcgmOne")
	entries["removedFunction"] = []string{"excluded", "removed"}
	problems := inventoryProblems(publicAPINames(header), entries, calls)
	if len(problems) != 2 {
		t.Fatalf("missing both-direction checks: %v", problems)
	}
}

type dynamicBridgeFixture struct {
	name, goSource, cSource string
	bound                   bool
}

// dynamicBridgeFixtures builds genuine and misleading bridge sources for inventory checks.
func dynamicBridgeFixtures(bridge dynamicAPIBridge) []dynamicBridgeFixture {
	parameters := "dcgmHandle_t, const char *, unsigned int *, unsigned int *"
	if bridge.symbol == "dcgmConfigReset" {
		parameters = "dcgmHandle_t, dcgmGpuGrp_t, dcgmStatus_t"
	}
	if bridge.symbol == "dcgmVersionInfo" {
		parameters = "dcgmVersionInfo_t *"
	}
	if bridge.symbol == "dcgmHostengineVersionInfo" {
		parameters = "dcgmHandle_t, dcgmVersionInfo_t *"
	}

	resolver := fmt.Sprintf(`%s resolved = (%s)dlsym(library, %q);`, bridge.functionType, bridge.functionType, bridge.symbol)
	cSource := fmt.Sprintf(`typedef dcgmReturn_t (*%s)(%s);
static dcgmReturn_t %s(void *library, dcgmHandle_t handle) {
 %s
 if (resolved == NULL) { return DCGM_ST_FUNCTION_NOT_FOUND; }
 return resolved(%s);
}`, bridge.functionType, parameters, bridge.bridge, resolver, bridge.arguments)
	goSource := "package dcgm\n/*\n" + cSource + "\n*/\nimport \"C\"\nfunc invoke() { C." + bridge.bridge + "() }"
	if bridge.cFile != bridge.goFile {
		goSource = "package dcgm\n/*\n#include \"" + bridge.cFile + "\"\n*/\nimport \"C\"\nfunc invoke() { C." + bridge.bridge + "() }"
	}

	invocation := "return resolved(" + bridge.arguments + ");"
	fixtures := []dynamicBridgeFixture{
		{"typed bridge", goSource, cSource, true},
		{"missing resolver", goSource, strings.ReplaceAll(cSource, "dlsym(", "otherResolver("), false},
		{"wrong symbol", goSource, strings.ReplaceAll(cSource, `"`+bridge.symbol+`"`, `"differentSymbol"`), false},
		{"missing typed invocation", goSource, strings.ReplaceAll(cSource, "return resolved(", "return unused("), false},
		{"wrong function type", goSource, strings.ReplaceAll(cSource, "typedef dcgmReturn_t", "typedef void"), false},
		{"missing bridge definition", goSource, strings.ReplaceAll(cSource, bridge.bridge+"(", "unusedBridge("), false},
		{"resolver outside bridge", goSource, resolver + "\n" + strings.Replace(cSource, resolver, "", 1), false},
		{"missing Go call", strings.ReplaceAll(goSource, "C."+bridge.bridge+"()", ""), cSource, false},
		{"commented Go call", strings.ReplaceAll(goSource, "C."+bridge.bridge+"()", "/* C."+bridge.bridge+"() */"), cSource, false},
		{"string is not Go call", strings.ReplaceAll(goSource, "C."+bridge.bridge+"()", `_ = "C.`+bridge.bridge+`()"`), cSource, false},
		{"unreachable direct call is not bridge evidence", strings.ReplaceAll(goSource, "C."+bridge.bridge+"()", "if false { C."+bridge.symbol+"() }"), cSource, false},
		{"body string with unmatched brace", goSource, strings.Replace(cSource, resolver, resolver+` const char *text = "}";`, 1), true},
		{"body string with comment markers", goSource, strings.Replace(cSource, invocation, `const char *text = "http://example/{/*"; `+invocation, 1), true},
		{"body char with unmatched brace", goSource, strings.Replace(cSource, resolver, resolver+` char brace = '}';`, 1), true},
		{"quoted invocation", goSource, strings.Replace(cSource, invocation, `const char *text = "`+invocation+`"; return DCGM_ST_OK;`, 1), false},
		{"quoted resolver", goSource, strings.Replace(cSource, resolver, `const char *text = '`+resolver+`';`, 1), false},
		{"wrong native arguments", goSource, strings.Replace(cSource, invocation, "return resolved();", 1), false},
	}
	for i := range fixtures {
		if bridge.goFile == bridge.cFile {
			fixtures[i].goSource = strings.Replace(fixtures[i].goSource, cSource, fixtures[i].cSource, 1)
		}
	}

	return fixtures
}

// checkDynamicBridgeFixture checks bridge evidence through the same source map as the inventory.
func checkDynamicBridgeFixture(t *testing.T, bridge dynamicAPIBridge, fixture dynamicBridgeFixture) {
	t.Helper()
	sources := map[string]string{bridge.goFile: fixture.goSource}
	if bridge.cFile != bridge.goFile {
		sources[bridge.cFile] = fixture.cSource
	}

	calls, err := inventoryNativeCalls(sources)
	if err != nil {
		t.Fatal(err)
	}

	entries := map[string][]string{bridge.symbol: {"bound", ""}}
	problems := inventoryProblems([]string{bridge.symbol}, entries, calls)
	if bound := len(problems) == 0; bound != fixture.bound {
		t.Fatalf("bound=%t, want %t: %v", bound, fixture.bound, problems)
	}
}

// TestPublicAPIDynamicBridgeEvidence rejects missing calls and decoys in optional native bridges.
func TestPublicAPIDynamicBridgeEvidence(t *testing.T) {
	for _, bridge := range dynamicAPIBridges {
		t.Run(bridge.symbol, func(t *testing.T) {
			for _, fixture := range dynamicBridgeFixtures(bridge) {
				t.Run(fixture.name, func(t *testing.T) { checkDynamicBridgeFixture(t, bridge, fixture) })
			}
		})
	}
}

// TestPublicAPIPIDBridgeEvidence requires the typed PID bridge rather than a textual call lookalike.
func TestPublicAPIPIDBridgeEvidence(t *testing.T) {
	bridge := `static inline dcgmReturn_t getPidInfoV2(dcgmHandle_t dcgmHandle, dcgmGpuGrp_t groupId, dcgmPidInfo_v2 *pidInfo) {
 return dcgmGetPidInfo(dcgmHandle, groupId, (dcgmPidInfo_t *)pidInfo);
}`
	goSource := "package dcgm\n/*\n" + bridge + "\n*/\nimport \"C\"\nfunc invoke() { C.getPidInfoV2() }"
	for _, fixture := range []struct {
		name, source string
		bound        bool
	}{
		{"typed bridge", goSource, true},
		{"missing Go call", strings.ReplaceAll(goSource, "C.getPidInfoV2()", ""), false},
		{"Go string is not a call", strings.ReplaceAll(goSource, "C.getPidInfoV2()", `_ = "C.getPidInfoV2()"`), false},
		{"missing definition", strings.Replace(goSource, bridge, "", 1), false},
		{"wrong argument type", strings.Replace(goSource, "dcgmPidInfo_v2 *", "dcgmPidInfo_v3 *", 1), false},
		{"wrong native call", strings.Replace(goSource, "return dcgmGetPidInfo(", "return otherNativeCall(", 1), false},
		{"missing native call", strings.Replace(goSource, "return dcgmGetPidInfo(dcgmHandle, groupId, (dcgmPidInfo_t *)pidInfo);", "return DCGM_ST_OK;", 1), false},
		{"missing cast", strings.Replace(goSource, "(dcgmPidInfo_t *)pidInfo", "pidInfo", 1), false},
		{"comment is not definition", strings.Replace(goSource, bridge, "// "+strings.ReplaceAll(bridge, "\n", "\n// "), 1), false},
		{"C string is not definition", strings.Replace(goSource, bridge, `const char *text = "`+strings.ReplaceAll(bridge, "\n", " ")+`";`, 1), false},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			calls, err := inventoryNativeCalls(map[string]string{"process_info.go": fixture.source})
			if err != nil {
				t.Fatal(err)
			}
			if calls["dcgmGetPidInfo"] != fixture.bound {
				t.Fatalf("PID bridge bound=%t, want %t", calls["dcgmGetPidInfo"], fixture.bound)
			}
		})
	}
}

// TestPublicAPINonCgoCallEvidence rejects package variables named C from other files.
func TestPublicAPINonCgoCallEvidence(t *testing.T) {
	sources := map[string]string{
		"declaration.go": "package dcgm\nvar C = struct{ dcgmOne func() }{func() {}}",
		"call.go":        "package dcgm\nfunc invoke() { C.dcgmOne() }",
	}
	calls, err := inventoryNativeCalls(sources)
	if err != nil {
		t.Fatal(err)
	}
	if calls["dcgmOne"] {
		t.Fatal("non-cgo file contributed native-call evidence through a package variable named C")
	}
}

// TestPublicAPIDirectCallEvidence counts only actual C calls in the Go syntax tree.
func TestPublicAPIDirectCallEvidence(t *testing.T) {
	for _, fixture := range []struct {
		name, source string
		bound        bool
	}{
		{"direct call", "C.dcgmOne()", true},
		{"local receiver", "C := struct{ dcgmOne func() }{func() {}}; C.dcgmOne()", false},
		{"double quoted string", `_ = "C.dcgmOne()"`, false},
		{"raw string", "_ = `C.dcgmOne()`", false},
		{"block comment", "/* C.dcgmOne() */", false},
		{"line comment", "// C.dcgmOne()\n", false},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			source := "package dcgm\nimport \"C\"\nfunc invoke() { " + fixture.source + " }"
			calls, err := inventoryNativeCalls(map[string]string{"direct.go": source, "unused.c": "C.dcgmOther()"})
			if err != nil {
				t.Fatal(err)
			}
			if calls["dcgmOne"] != fixture.bound || calls["dcgmOther"] {
				t.Fatalf("native calls %v, want dcgmOne=%t only", calls, fixture.bound)
			}
		})
	}
}
