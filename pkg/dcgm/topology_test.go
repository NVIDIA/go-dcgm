package dcgm

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestDeviceTopologyFromNative checks empty, populated, and full responses and verifies that the links are Go-owned.
func TestDeviceTopologyFromNative(t *testing.T) {
	capacity := len(deviceTopologyResponseForTest(0).gpuPaths)
	for _, count := range []int{0, 2, capacity} {
		t.Run(fmt.Sprintf("count_%d", count), func(t *testing.T) {
			response := deviceTopologyResponseForTest(uint(count))
			if count > 0 {
				response.gpuPaths[0].gpuId = 7
				response.gpuPaths[0].path = 0x200 // Two NVLinks.
				response.gpuPaths[count-1].gpuId = 29
			}

			got, err := deviceTopologyFromNative(response)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got, count)
			if count > 0 {
				require.Equal(t, P2PLink{GPU: 7, Link: TwoNVLINKLinks}, got[0])
				require.Equal(t, P2PLink{GPU: 29, Link: P2PLinkUnknown}, got[count-1])
			}

			want := slices.Clone(got)
			*response = *deviceTopologyResponseForTest(0)
			require.Equal(t, want, got)
		})
	}
}

// TestDeviceTopologyFromNativeRejectsInvalidCounts checks that an oversized peer count returns an error without links.
func TestDeviceTopologyFromNativeRejectsInvalidCounts(t *testing.T) {
	capacity := len(deviceTopologyResponseForTest(0).gpuPaths)
	for _, count := range []uint{uint(capacity + 1), uint(^uint32(0))} {
		t.Run(fmt.Sprintf("count_%d", count), func(t *testing.T) {
			got, err := deviceTopologyFromNative(deviceTopologyResponseForTest(count))
			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf("invalid topology peer count %d, expected 0 through %d", count, capacity))
		})
	}
}

// TestNvLinkStatusFromNative checks link order, array limits, and ownership of the returned statuses.
func TestNvLinkStatusFromNative(t *testing.T) {
	empty := nvLinkStatusResponseForTest(0, 0)
	gpuCapacity := len(empty.gpus)
	switchCapacity := len(empty.nvSwitches)
	gpuLinks := len(empty.gpus[0].linkState)
	switchLinks := len(empty.nvSwitches[0].linkState)

	for _, tt := range []struct {
		name     string
		gpus     int
		switches int
	}{
		{name: "empty"},
		{name: "mixed", gpus: 2, switches: 2},
		{name: "GPU capacity", gpus: gpuCapacity},
		{name: "switch capacity", switches: switchCapacity},
		{name: "both capacities", gpus: gpuCapacity, switches: switchCapacity},
	} {
		t.Run(tt.name, func(t *testing.T) {
			response := nvLinkStatusResponseForTest(uint(tt.gpus), uint(tt.switches))
			response.gpus[0].entityId = 7
			response.gpus[0].linkState[0] = nativeNVLinkStateUp
			response.nvSwitches[0].entityId = 19
			response.nvSwitches[0].linkState[0] = nativeNVLinkStateDown
			if tt.gpus > 0 {
				response.gpus[tt.gpus-1].entityId = 29
				response.gpus[tt.gpus-1].linkState[gpuLinks-1] = nativeNVLinkStateDisabled
			}
			if tt.switches > 0 {
				response.nvSwitches[tt.switches-1].entityId = 31
				response.nvSwitches[tt.switches-1].linkState[switchLinks-1] = nativeNVLinkStateUp
			}

			got, err := nvLinkStatusFromNative(response)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Len(t, got, tt.gpus*gpuLinks+tt.switches*switchLinks)
			if tt.gpus > 0 {
				require.Equal(t, NvLinkStatus{ParentId: 7, ParentType: FE_GPU, State: LS_UP}, got[0])
				require.Equal(t, NvLinkStatus{
					ParentId: 29, ParentType: FE_GPU, State: LS_DISABLED, Index: uint(gpuLinks - 1),
				}, got[tt.gpus*gpuLinks-1])
			}
			if tt.switches > 0 {
				require.Equal(t, NvLinkStatus{ParentId: 19, ParentType: FE_SWITCH, State: LS_DOWN}, got[tt.gpus*gpuLinks])
				require.Equal(t, NvLinkStatus{
					ParentId: 31, ParentType: FE_SWITCH, State: LS_UP, Index: uint(switchLinks - 1),
				}, got[len(got)-1])
			}

			want := slices.Clone(got)
			*response = *empty
			require.Equal(t, want, got)
		})
	}
}

// TestNvLinkStatusFromNativeRejectsInvalidCounts checks that either oversized count rejects the entire response.
func TestNvLinkStatusFromNativeRejectsInvalidCounts(t *testing.T) {
	gpuCapacity := len(nvLinkStatusResponseForTest(0, 0).gpus)
	switchCapacity := len(nvLinkStatusResponseForTest(0, 0).nvSwitches)
	oversizedGPUCount := uint(gpuCapacity + 1)
	oversizedSwitchCount := uint(switchCapacity + 1)
	const maxCount = uint(^uint32(0))

	for _, tt := range []struct {
		name     string
		gpus     uint
		switches uint
		capacity int
		count    uint
		kind     string
	}{
		{
			name:     "GPU capacity plus one",
			gpus:     oversizedGPUCount,
			switches: 1,
			capacity: gpuCapacity,
			count:    oversizedGPUCount,
			kind:     "GPU",
		},
		{
			name:     "GPU maximum count",
			gpus:     maxCount,
			switches: 1,
			capacity: gpuCapacity,
			count:    maxCount,
			kind:     "GPU",
		},
		{
			name:     "switch capacity plus one",
			gpus:     1,
			switches: oversizedSwitchCount,
			capacity: switchCapacity,
			count:    oversizedSwitchCount,
			kind:     "switch",
		},
		{
			name:     "switch maximum count",
			gpus:     1,
			switches: maxCount,
			capacity: switchCapacity,
			count:    maxCount,
			kind:     "switch",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nvLinkStatusFromNative(nvLinkStatusResponseForTest(tt.gpus, tt.switches))
			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf("invalid NVLink %s count %d, expected 0 through %d", tt.kind, tt.count, tt.capacity))
		})
	}
}

// TestGetDeviceTopologyWithOps checks peer GPU bus ID lookup and errors from DCGM reads.
func TestGetDeviceTopologyWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")

	tests := []struct {
		name         string
		links        []P2PLink
		readErr      error
		lookup       bool
		wantEntities []GroupEntityPair
		values       []FieldValue_v2
		lookupErr    error
		want         []P2PLink
		wantErrIs    error
		wantErrText  string
	}{
		{
			name:      "returns native read error",
			readErr:   stageErr,
			wantErrIs: stageErr,
		},
		{
			name: "skips lookup for nil links",
		},
		{
			name:  "skips lookup for empty links",
			links: []P2PLink{},
			want:  []P2PLink{},
		},
		{
			name:   "wraps peer lookup error",
			links:  []P2PLink{{GPU: 2, Link: P2PLinkSingleSwitch}, {GPU: 7, Link: TwoNVLINKLinks}},
			lookup: true,
			wantEntities: []GroupEntityPair{
				{EntityGroupId: FE_GPU, EntityId: 2},
				{EntityGroupId: FE_GPU, EntityId: 7},
			},
			lookupErr:   stageErr,
			wantErrIs:   stageErr,
			wantErrText: "get peer GPU bus IDs",
		},
		{
			name:         "returns missing peer error",
			links:        []P2PLink{{GPU: 2, Link: P2PLinkSingleSwitch}},
			lookup:       true,
			wantEntities: []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: 2}},
			wantErrText:  "get bus ID for peer GPU 2: not returned by DCGM",
		},
		{
			name:   "enriches every peer from out-of-order values",
			links:  []P2PLink{{GPU: 2, Link: P2PLinkSingleSwitch}, {GPU: 7, Link: TwoNVLINKLinks}},
			lookup: true,
			wantEntities: []GroupEntityPair{
				{EntityGroupId: FE_GPU, EntityId: 2},
				{EntityGroupId: FE_GPU, EntityId: 7},
			},
			values: []FieldValue_v2{
				fieldValueV2String(FE_GPU, 7, DCGM_FI_DEV_PCI_BUS_ID, "0000:07:00.0"),
				fieldValueV2String(FE_GPU, 2, DCGM_FI_DEV_PCI_BUS_ID, "0000:02:00.0"),
			},
			want: []P2PLink{
				{GPU: 2, BusID: "0000:02:00.0", Link: P2PLinkSingleSwitch},
				{GPU: 7, BusID: "0000:07:00.0", Link: TwoNVLINKLinks},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewMockTopologyOps(gomock.NewController(t))
			api.EXPECT().readDeviceTopology(uint(3)).Return(tt.links, tt.readErr)
			if tt.lookup {
				api.EXPECT().entitiesGetLatestValues(
					tt.wantEntities,
					[]Short{DCGM_FI_DEV_PCI_BUS_ID},
					DCGM_FV_FLAG_LIVE_DATA,
				).Return(tt.values, tt.lookupErr)
			}

			got, err := getDeviceTopologyWithOps(api, 3)
			if tt.wantErrText != "" || tt.wantErrIs != nil {
				require.Error(t, err)
				require.Nil(t, got)
				if tt.wantErrText != "" {
					require.ErrorContains(t, err, tt.wantErrText)
				}
				if tt.wantErrIs != nil {
					require.ErrorIs(t, err, tt.wantErrIs)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// TestGetDeviceTopologyRejectsInvalidPeerFields checks that a bad peer field
// fails the whole topology read and retains any native field status.
func TestGetDeviceTopologyRejectsInvalidPeerFields(t *testing.T) {
	for _, test := range []struct {
		name      string
		fieldType uint
		status    int
		want      string
	}{
		{name: "failed required field", fieldType: DCGM_FT_STRING, status: DCGM_ST_NO_DATA, want: "field status"},
		{name: "wrong field type", fieldType: DCGM_FT_INT64, want: "expected a string value"},
	} {
		t.Run(test.name, func(t *testing.T) {
			const sourceGPU = 3
			const invalidPeer = 7
			const busIDField = DCGM_FI_DEV_PCI_BUS_ID

			// The first peer has valid data so the test also checks that no partial result escapes.
			links := []P2PLink{{GPU: 2}, {GPU: invalidPeer}}
			gpus := []GroupEntityPair{
				{EntityGroupId: FE_GPU, EntityId: 2},
				{EntityGroupId: FE_GPU, EntityId: invalidPeer},
			}
			valid := fieldValueV2String(FE_GPU, 2, busIDField, "0000:02:00.0")
			invalid := fieldValueV2String(FE_GPU, invalidPeer, busIDField, "0000:07:00.0")
			invalid.FieldType, invalid.Status = test.fieldType, test.status

			api := NewMockTopologyOps(gomock.NewController(t))
			api.EXPECT().readDeviceTopology(uint(sourceGPU)).Return(links, nil)
			api.EXPECT().entitiesGetLatestValues(gpus, []Short{DCGM_FI_DEV_PCI_BUS_ID}, gomock.Any()).
				Return([]FieldValue_v2{valid, invalid}, nil)
			got, err := getDeviceTopologyWithOps(api, sourceGPU)

			require.Nil(t, got)
			context := fmt.Sprintf("get topology for GPU %d: get bus ID for peer GPU %d, field %d",
				sourceGPU, invalidPeer, busIDField)
			require.ErrorContains(t, err, context)
			require.ErrorContains(t, err, test.want)
			if test.status != DCGM_ST_OK {
				var native *Error
				require.ErrorAs(t, err, &native)
				require.Equal(t, test.status, int(native.Code))
			}
		})
	}
}

func TestSelectedPopulatePeerBusIDs(t *testing.T) {
	tests := []struct {
		name        string
		links       []P2PLink
		values      []FieldValue_v2
		wantBusIDs  map[uint]string
		wantErrText string
	}{
		{
			name:  "maps each peer GPU from a batch response",
			links: []P2PLink{{GPU: 2}, {GPU: 7}},
			values: []FieldValue_v2{
				fieldValueV2String(FE_GPU, 7, DCGM_FI_DEV_PCI_BUS_ID, "0000:07:00.0"),
				fieldValueV2String(FE_GPU, 2, DCGM_FI_DEV_PCI_BUS_ID, "0000:02:00.0"),
			},
			wantBusIDs: map[uint]string{2: "0000:02:00.0", 7: "0000:07:00.0"},
		},
		{
			name:        "reports a missing peer GPU",
			links:       []P2PLink{{GPU: 2}},
			wantErrText: "peer GPU 2",
		},
		{
			name:  "reports a non-OK field status",
			links: []P2PLink{{GPU: 2}},
			values: []FieldValue_v2{
				{
					EntityGroupId: FE_GPU,
					EntityID:      2,
					FieldID:       DCGM_FI_DEV_PCI_BUS_ID,
					FieldType:     DCGM_FT_STRING,
					Status:        1,
				},
			},
			wantErrText: "field status 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := populatePeerBusIDs(tt.links, tt.values)
			if tt.wantErrText != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("populatePeerBusIDs() error = %v, want text %q", err, tt.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("populatePeerBusIDs() error = %v", err)
			}

			for i := range tt.links {
				if got, want := tt.links[i].BusID, tt.wantBusIDs[tt.links[i].GPU]; got != want {
					t.Errorf("links[%d].BusID = %q, want %q", i, got, want)
				}
			}
		})
	}
}

func fieldValueV2String(group Field_Entity_Group, entityID uint, fieldID Short, value string) FieldValue_v2 {
	fieldValue := FieldValue_v2{
		EntityGroupId: group,
		EntityID:      entityID,
		FieldID:       fieldID,
		FieldType:     DCGM_FT_STRING,
		Status:        DCGM_ST_OK,
	}
	copy(fieldValue.Value[:], value)
	return fieldValue
}

func TestPeerBusIDRejectsUnavailableStrings(t *testing.T) {
	for _, test := range []struct {
		name string
		text string
	}{
		{name: "empty"},
		{name: "spaces", text: "   "},
		{name: "tabs and newlines", text: "\t\n"},
		{name: "blank sentinel", text: DCGM_FT_STR_BLANK},
		{name: "not found sentinel", text: DCGM_FT_STR_NOT_FOUND},
		{name: "not supported sentinel", text: DCGM_FT_STR_NOT_SUPPORTED},
		{name: "permission sentinel", text: "<<<NOT_PERM>>>"},
		{name: "other DCGM sentinel", text: "<<<UNAVAILABLE>>>"},
	} {
		t.Run(test.name, func(t *testing.T) {
			value := fieldValueV2String(FE_GPU, 7, DCGM_FI_DEV_PCI_BUS_ID, test.text)
			require.ErrorContains(t, populatePeerBusIDs([]P2PLink{{GPU: 7}}, []FieldValue_v2{value}), "empty or blank value")
		})
	}
	t.Run("NUL outside the native string member", func(t *testing.T) {
		valueWithoutTerminator := strings.Repeat("x", fieldStringCapacity)
		value := fieldValueV2String(FE_GPU, 7, DCGM_FI_DEV_PCI_BUS_ID, valueWithoutTerminator)
		links := []P2PLink{{GPU: 7}}
		require.ErrorContains(t, populatePeerBusIDs(links, []FieldValue_v2{value}), "unterminated string")
		require.Empty(t, links[0].BusID)

		copy(value.Value[fieldStringCapacity:], "ignored union bytes")
		require.Equal(t, valueWithoutTerminator, value.String())
	})
}

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
		`(?s)func \(cgoAdapter\) readDeviceTopology\(gpuID uint\).*var topology C\.dcgmDeviceTopology_v2.*topology\.version = makeVersion2\(unsafe\.Sizeof\(topology\)\)`)
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
