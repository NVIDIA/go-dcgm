package dcgm

import (
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestLatestValuesForDeviceReportsCleanupFailures(t *testing.T) {
	collectionErr := errors.New("get latest values")
	for _, tt := range []struct {
		name        string
		getErr      error
		fieldStatus int
		want        DeviceStatus
	}{
		{name: "successful collection", want: DeviceStatus{Temperature: 42}},
		{name: "collection failure", getErr: collectionErr},
		{
			name:        "partial collection",
			fieldStatus: DCGM_ST_GPU_IS_LOST,
			want:        DeviceStatus{Temperature: 42, Power: DCGM_FT_FP64_BLANK},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fieldErr := errors.New("destroy fields")
			groupErr := errors.New("destroy group")
			m := NewMockDeviceStatusOps(gomock.NewController(t))
			create := m.EXPECT().fieldGroupCreate(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, fields []Short) (FieldHandle, error) {
					if len(fields) != 17 {
						t.Fatalf("field count = %d, want 17", len(fields))
					}
					return FieldHandle{}, nil
				},
			)
			watch := m.EXPECT().watchFields(uint(3), FieldHandle{}, gomock.Any()).Return(GroupHandle{}, nil).After(create)
			collect := m.EXPECT().getLatestValuesForFields(uint(3), gomock.Any()).DoAndReturn(
				func(_ uint, fields []Short) ([]FieldValue_v1, error) {
					if tt.getErr != nil {
						return nil, tt.getErr
					}
					values := deviceStatusSamples()
					binary.NativeEndian.PutUint64(values[deviceTemperature].Value[:], 42)
					values[devicePower].Status = tt.fieldStatus
					return values, nil
				},
			).After(watch)
			m.EXPECT().fieldGroupDestroy(FieldHandle{}).Return(fieldErr).After(collect)
			m.EXPECT().destroyGroup(GroupHandle{}).Return(groupErr).After(collect)
			status, err := latestValuesForDeviceWithOps(3, m)
			if status != tt.want {
				t.Fatalf("status = %+v, want %+v", status, tt.want)
			}
			if tt.fieldStatus != DCGM_ST_OK {
				nativeErr, ok := errors.AsType[*Error](err)
				if !ok || int(nativeErr.Code) != tt.fieldStatus {
					t.Errorf("error = %v, want DCGM status %d", err, tt.fieldStatus)
				}
				if err == nil || !strings.Contains(err.Error(), "get status for GPU 3") {
					t.Errorf("error = %v, want GPU context", err)
				}
			}
			for _, want := range []error{tt.getErr, fieldErr, groupErr} {
				if want != nil && !errors.Is(err, want) {
					t.Errorf("error = %v, missing %v", err, want)
				}
			}
		})
	}
}
