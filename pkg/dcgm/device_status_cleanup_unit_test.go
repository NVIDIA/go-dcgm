package dcgm

import (
	"encoding/binary"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
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
					values := make([]FieldValue_v1, len(fields))
					binary.NativeEndian.PutUint64(values[1].Value[:], 42)
					return values, nil
				},
			).After(watch)
			m.EXPECT().fieldGroupDestroy(FieldHandle{}).Return(fieldErr).After(collect)
			m.EXPECT().destroyGroup(GroupHandle{}).Return(groupErr).After(collect)
			status, err := latestValuesForDeviceWithOps(3, m)
			if status != tt.want {
				t.Fatalf("status = %+v, want %+v", status, tt.want)
			}
			for _, want := range []error{tt.getErr, fieldErr, groupErr} {
				if want != nil && !errors.Is(err, want) {
					t.Errorf("error = %v, missing %v", err, want)
				}
			}
		})
	}
}
