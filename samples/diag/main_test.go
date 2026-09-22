package main

import (
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func TestDiagOutputUsesCurrentResults(t *testing.T) {
	results := dcgm.DiagResults{
		DCGMVersion:   "4.6.1",
		DriverVersion: "580.1",
		Entities: []dcgm.DiagEntity{{
			Entity:       dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_GPU, EntityId: 0},
			SerialNumber: "GPU-SERIAL",
			SKUDeviceID:  "27B8",
		}},
		SystemErrors: []dcgm.DiagError{{
			Entity:  dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_NONE},
			Code:    99,
			Message: "system failure",
		}},
		Tests: []dcgm.DiagTest{{
			Name:       "memory",
			PluginName: "memory",
			Category:   "Hardware",
			Status:     "pass",
			Results: []dcgm.DiagEntityResult{{
				Entity: dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_GPU, EntityId: 0},
				Status: "pass",
			}},
			Errors: []dcgm.DiagError{{
				Entity:   dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_NONE},
				Code:     42,
				Category: dcgm.DCGM_FR_EC_SOFTWARE_CONFIG,
				Severity: dcgm.DCGM_ERROR_CONFIG,
				Message:  "warning",
			}},
			Info: []dcgm.DiagInfo{{
				Entity:  dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_NONE},
				Message: "global info",
			}, {
				Entity:  dcgm.GroupEntityPair{EntityGroupId: dcgm.FE_GPU, EntityId: 0},
				Message: "memory ok",
			}},
			AuxDataVersion: 1,
			AuxData:        `{"version":"1"}`,
		}},
	}

	var output strings.Builder
	if err := newDiagOutputTemplate().Execute(&output, results); err != nil {
		t.Fatalf("execute diagnostic output template: %v", err)
	}
	want := "DCGM 4.6.1\tDriver 580.1\n" +
		"Entity GPU 0\tSerial GPU-SERIAL\tSKU 27B8\n" +
		"System error\tcode=99\tcategory=0\tseverity=0\tsystem failure\n" +
		"[Hardware] memory\tpass\tPlugin memory\n" +
		"  Result GPU 0\tpass\n" +
		"  Error global\tcode=42\tcategory=3\tseverity=5\twarning\n" +
		"  Info global\tglobal info\n" +
		"  Info GPU 0\tmemory ok\n" +
		"  Auxiliary data\tversion=1\t{\"version\":\"1\"}\n"
	if got := output.String(); got != want {
		t.Fatalf("unexpected diagnostic output:\n%s\nwant:\n%s", got, want)
	}
}
