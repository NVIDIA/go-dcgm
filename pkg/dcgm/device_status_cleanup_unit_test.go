package dcgm

import (
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

func TestLatestValuesForDeviceReportsCleanupFailures(t *testing.T) {
	collectionErr := errors.New("get latest values")
	for _, tt := range []struct {
		name   string
		getErr error
		want   DeviceStatus
	}{
		{name: "successful collection", want: DeviceStatus{Temperature: 42}},
		{name: "collection failure", getErr: collectionErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fieldErr := errors.New("destroy fields")
			groupErr := errors.New("destroy group")
			var calls []string
			status, err := latestValuesForDeviceWithOps(3, deviceStatusOps{
				createFieldGroup: func(_ string, fields []Short) (FieldHandle, error) {
					calls = append(calls, "create fields")
					if len(fields) != 17 {
						t.Fatalf("field count = %d, want 17", len(fields))
					}
					return FieldHandle{}, nil
				},
				watchFields: func(gpuID uint, _ FieldHandle, _ string) (GroupHandle, error) {
					calls = append(calls, "watch fields")
					if gpuID != 3 {
						t.Fatalf("GPU ID = %d, want 3", gpuID)
					}
					return GroupHandle{}, nil
				},
				getLatestValues: func(_ uint, fields []Short) ([]FieldValue_v1, error) {
					calls = append(calls, "collect values")
					if tt.getErr != nil {
						return nil, tt.getErr
					}
					values := make([]FieldValue_v1, len(fields))
					binary.NativeEndian.PutUint64(values[1].Value[:], 42)
					return values, nil
				},
				destroyFieldGroup: func(FieldHandle) error {
					calls = append(calls, "destroy fields")
					return fieldErr
				},
				destroyGroup: func(GroupHandle) error {
					calls = append(calls, "destroy group")
					return groupErr
				},
			})
			if status != tt.want {
				t.Fatalf("status = %+v, want %+v", status, tt.want)
			}
			for _, want := range []error{tt.getErr, fieldErr, groupErr} {
				if want != nil && !errors.Is(err, want) {
					t.Errorf("error = %v, missing %v", err, want)
				}
			}
			wantCalls := []string{"create fields", "watch fields", "collect values", "destroy fields", "destroy group"}
			if !reflect.DeepEqual(calls, wantCalls) {
				t.Errorf("calls = %v, want %v", calls, wantCalls)
			}
		})
	}
}
