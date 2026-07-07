package dcgm

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestGetP2PLinkMapsTopologyConstants(t *testing.T) {
	constants := dcgmTopologyConstants(t)

	tests := []struct {
		name string
		want P2PLinkType
		path string
	}{
		{name: "DCGM_TOPOLOGY_BOARD", want: P2PLinkSameBoard, path: "PSB"},
		{name: "DCGM_TOPOLOGY_SINGLE", want: P2PLinkSingleSwitch, path: "PIX"},
		{name: "DCGM_TOPOLOGY_MULTIPLE", want: P2PLinkMultiSwitch, path: "PXB"},
		{name: "DCGM_TOPOLOGY_HOSTBRIDGE", want: P2PLinkHostBridge, path: "PHB"},
		{name: "DCGM_TOPOLOGY_CPU", want: P2PLinkSameCPU, path: "NODE"},
		{name: "DCGM_TOPOLOGY_SYSTEM", want: P2PLinkCrossCPU, path: "SYS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getP2PLink(requireTopologyConstant(t, constants, tt.name))
			if got != tt.want {
				t.Fatalf("getP2PLink(%s) = %v, want %v", tt.name, got, tt.want)
			}
			if gotPath := got.PCIPaths(); gotPath != tt.path {
				t.Fatalf("%s PCIPaths() = %q, want %q", tt.name, gotPath, tt.path)
			}
		})
	}

	for count := uint(1); count <= maxNVLinkCount; count++ {
		name := fmt.Sprintf("DCGM_TOPOLOGY_NVLINK%d", count)
		t.Run(name, func(t *testing.T) {
			got := getP2PLink(requireTopologyConstant(t, constants, name))
			wantPath := fmt.Sprintf("NV%d", count)
			if gotPath := got.PCIPaths(); gotPath != wantPath {
				t.Fatalf("%s PCIPaths() = %q, want %q", name, gotPath, wantPath)
			}
		})
	}

	t.Run("CompositeNVLinkAndPCIPathPrefersNVLink", func(t *testing.T) {
		path := requireTopologyConstant(t, constants, "DCGM_TOPOLOGY_CPU") |
			requireTopologyConstant(t, constants, "DCGM_TOPOLOGY_NVLINK2")
		if gotPath := getP2PLink(path).PCIPaths(); gotPath != "NV2" {
			t.Fatalf("composite topology path PCIPaths() = %q, want %q", gotPath, "NV2")
		}
	})

	t.Run("Unknown", func(t *testing.T) {
		tests := []uint64{
			requireTopologyConstant(t, constants, "DCGM_TOPOLOGY_UNINITIALIZED"),
			1 << 60,
		}
		for _, path := range tests {
			if got := getP2PLink(path); got != P2PLinkUnknown {
				t.Fatalf("getP2PLink(%#x) = %v, want %v", path, got, P2PLinkUnknown)
			}
		}
	})
}

func TestTopologyUsesLatestStructVersions(t *testing.T) {
	structsHeader := readTopologyTestFile(t, "dcgm_structs.h")
	assertContains(t, structsHeader, "typedef dcgmDeviceTopology_v2 dcgmDeviceTopology_t;")
	assertContains(t, structsHeader, "#define dcgmDeviceTopology_version dcgmDeviceTopology_version2")
	assertContains(t, structsHeader, "#define dcgmNvLinkStatus_version dcgmNvLinkStatus_version5")

	agentHeader := readTopologyTestFile(t, "dcgm_agent.h")
	assertContains(t, agentHeader, "dcgmNvLinkStatus_v5 *linkStatus")
	assertContains(t, agentHeader, "dcgmDeviceTopology_v2 *pDcgmDeviceTopology")

	internalHeader := readTopologyTestFile(t, "dcgm_structs_internal.h")
	assertContains(t, internalHeader, "DCGM_CASSERT(dcgmDeviceTopology_version2")
	assertContains(t, internalHeader, "DCGM_CASSERT(dcgmNvLinkStatus_version5")

	source := readTopologyTestFile(t, "topology.go")
	assertRegexp(t, source,
		`(?s)func getDeviceTopology\(gpuID uint\).*var topology C\.dcgmDeviceTopology_v2.*topology\.version = makeVersion2\(unsafe\.Sizeof\(topology\)\)`)
	assertRegexp(t, source,
		`(?s)func getNvLinkLinkStatus\(\).*var linkStatus C\.dcgmNvLinkStatus_v5.*linkStatus\.version = makeVersion5\(unsafe\.Sizeof\(linkStatus\)\)`)
}

func TestEntitiesGetLatestValuesV4StaysWithinProtocolLimit(t *testing.T) {
	internalSource := readTopologyTestFile(t, "internal.go")
	assertContains(t, internalSource, "dcgmEntitiesGetLatestValues_v4_exceeds_proto_limit")
}

func dcgmTopologyConstants(t *testing.T) map[string]uint64 {
	t.Helper()

	content := readTopologyTestFile(t, "dcgm_structs.h")
	constants := make(map[string]uint64)

	enumPattern := regexp.MustCompile(`(?m)^\s*(DCGM_TOPOLOGY_[A-Z0-9_]+)\s*=\s*(0x[0-9A-Fa-f]+|\d+),`)
	for _, match := range enumPattern.FindAllStringSubmatch(content, -1) {
		constants[match[1]] = parseTopologyConstantValue(t, match[1], match[2])
	}

	definePattern := regexp.MustCompile(`(?m)^#define\s+(DCGM_TOPOLOGY_[A-Z0-9_]+)\s+\(\(dcgmGpuTopologyLevel_t\)(0x[0-9A-Fa-f]+)ULL\)`)
	for _, match := range definePattern.FindAllStringSubmatch(content, -1) {
		constants[match[1]] = parseTopologyConstantValue(t, match[1], match[2])
	}

	if len(constants) == 0 {
		t.Fatal("no DCGM_TOPOLOGY_* constants found in dcgm_structs.h")
	}
	return constants
}

func parseTopologyConstantValue(t *testing.T, name, value string) uint64 {
	t.Helper()

	parsed, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		t.Fatalf("parse %s value %q: %v", name, value, err)
	}
	return parsed
}

func requireTopologyConstant(t *testing.T, constants map[string]uint64, name string) uint64 {
	t.Helper()

	value, ok := constants[name]
	if !ok {
		t.Fatalf("%s was not found in dcgm_structs.h", name)
	}
	return value
}

func readTopologyTestFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertContains(t *testing.T, content, needle string) {
	t.Helper()

	if !strings.Contains(content, needle) {
		t.Fatalf("expected file content to contain %q", needle)
	}
}

func assertRegexp(t *testing.T, content, pattern string) {
	t.Helper()

	if !regexp.MustCompile(pattern).MatchString(content) {
		t.Fatalf("expected file content to match %q", pattern)
	}
}
