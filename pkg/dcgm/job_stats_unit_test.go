package dcgm

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestJobIDValidation(t *testing.T) {
	for _, id := range []string{"", strings.Repeat("x", 64), "a\x00b", strings.Repeat("é", 32)} {
		t.Run(id, func(t *testing.T) {
			api := NewMockJobStatsCalls(gomock.NewController(t))
			require.Error(t, jobStartWithOps(api, GroupHandle{}, id))
			require.Error(t, jobStopWithOps(api, id))
			_, err := jobGetStatsWithOps(api, id)
			require.Error(t, err)
			require.Error(t, jobRemoveWithOps(api, id))
		})
	}
	for _, id := range []string{"x", strings.Repeat("x", 63), strings.Repeat("é", 31)} {
		t.Run(id, func(t *testing.T) {
			api := NewMockJobStatsCalls(gomock.NewController(t))
			api.EXPECT().jobStartNative(GroupHandle{}, id)
			api.EXPECT().jobStopNative(id)
			api.EXPECT().jobGetStatsNative(id).Return(JobInfo{GPUs: []GPUJobStats{{GPU: 7}}}, nil)
			api.EXPECT().jobRemoveNative(id)
			require.NoError(t, jobStartWithOps(api, GroupHandle{}, id))
			require.NoError(t, jobStopWithOps(api, id))
			info, err := jobGetStatsWithOps(api, id)
			require.NoError(t, err)
			require.Equal(t, uint(7), info.GPUs[0].GPU)
			require.NoError(t, jobRemoveWithOps(api, id))
		})
	}
}

func TestJobErrors(t *testing.T) {
	api := NewMockJobStatsCalls(gomock.NewController(t))
	want := jobNativeErrorForTest()
	var dcgmErr *Error
	require.ErrorAs(t, want, &dcgmErr)
	require.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
	api.EXPECT().watchJobFieldsNative(GroupHandle{}, time.Second, time.Minute, 10).Return(want)
	api.EXPECT().jobStartNative(GroupHandle{}, "job").Return(want)
	api.EXPECT().jobStopNative("job").Return(want)
	api.EXPECT().jobGetStatsNative("job").Return(JobInfo{}, want)
	api.EXPECT().jobRemoveNative("job").Return(want)
	require.ErrorIs(t, watchJobFieldsWithOps(api, GroupHandle{}, time.Second, time.Minute, 10), want)
	require.ErrorIs(t, jobStartWithOps(api, GroupHandle{}, "job"), want)
	require.ErrorIs(t, jobStopWithOps(api, "job"), want)
	_, err := jobGetStatsWithOps(api, "job")
	require.ErrorIs(t, err, want)
	require.ErrorIs(t, jobRemoveWithOps(api, "job"), want)
}

func TestJobWatchValidation(t *testing.T) {
	for _, tc := range []struct {
		freq, age time.Duration
		samples   int
	}{
		{0, 0, 0},
		{-time.Second, 0, 0},
		{time.Nanosecond, 0, 0},
		{1500 * time.Nanosecond, 0, 0},
		{time.Second + time.Nanosecond, 0, 0},
		{time.Second, -time.Second, 0},
		{time.Second, 0, -1},
		{time.Second, 0, math.MaxInt32 + 1},
	} {
		// A mock with no expectations rejects any native call for invalid input.
		require.ErrorContains(t, watchJobFieldsWithOps(NewMockJobStatsCalls(gomock.NewController(t)), GroupHandle{}, tc.freq, tc.age, tc.samples), "watch job fields:")
	}
	api := NewMockJobStatsCalls(gomock.NewController(t))
	api.EXPECT().watchJobFieldsNative(GroupHandle{}, time.Microsecond, time.Duration(0), 0)
	require.NoError(t, watchJobFieldsWithOps(api, GroupHandle{}, time.Microsecond, 0, 0))
	api.EXPECT().watchJobFieldsNative(GroupHandle{}, time.Second, 1500*time.Millisecond, 10)
	require.NoError(t, watchJobFieldsWithOps(api, GroupHandle{}, time.Second, 1500*time.Millisecond, 10))
	maximumAge := time.Duration(math.MaxInt64)
	require.False(t, math.IsInf(maximumAge.Seconds(), 0))
	api.EXPECT().watchJobFieldsNative(GroupHandle{}, time.Second, maximumAge, math.MaxInt32)
	require.NoError(t, watchJobFieldsWithOps(api, GroupHandle{}, time.Second, maximumAge, math.MaxInt32))
}

func TestJobConversion(t *testing.T) {
	info, err := jobResponseForTest(1, 1, 1, 1, 1)
	require.NoError(t, err)
	require.Len(t, info.GPUs, 1)
	g := info.GPUs[0]
	require.Equal(t, uint(107), info.Summary.GPU)
	require.Equal(t, int64(131), info.Summary.PCIeReplays)
	require.Equal(t, int64(1000101), info.Summary.StartTime)
	require.Equal(t, int64(2000102), info.Summary.EndTime)
	require.Equal(t, uint(103), info.Summary.ECCSingleBit)
	require.Equal(t, uint(104), info.Summary.ECCDoubleBit)
	require.Equal(t, int64(987654421), info.Summary.MaxGPUMemoryUsed)
	require.Equal(t, []int64{201, 202, 203, 204, 205, 206}, []int64{info.Summary.PowerViolationTime, info.Summary.ThermalViolationTime, info.Summary.ReliabilityViolationTime, info.Summary.BoardLimitViolationTime, info.Summary.LowUtilizationTime, info.Summary.SyncBoostTime})
	require.Equal(t, int64(12445), info.Summary.EnergyConsumed)
	require.Equal(t, JobFloat64Summary{101.25, 109.75, 104.5}, info.Summary.PowerUsage)
	require.Equal(t, JobInt64Summary{111, 122, 117}, info.Summary.PCIeRxBandwidth)
	require.Equal(t, JobInt64Summary{141, 152, 147}, info.Summary.PCIeTxBandwidth)
	require.Equal(t, JobInt32Summary{102, 195, 160}, info.Summary.SMUtilization)
	require.Equal(t, JobInt32Summary{103, 175, 140}, info.Summary.MemoryUtilization)
	require.Equal(t, JobInt32Summary{1100, 2100, 1600}, info.Summary.MemoryClock)
	require.Equal(t, JobInt32Summary{500, 1300, 900}, info.Summary.SMClock)
	require.Equal(t, []JobProcessStats{{223, 103.5, 104.5}}, info.Summary.ComputeProcesses)
	require.Equal(t, []JobProcessStats{{556, 107.5, 108.5}}, info.Summary.GraphicsProcesses)
	require.Equal(t, uint(7), g.GPU)
	require.Equal(t, int64(12345), g.EnergyConsumed)
	require.Equal(t, JobFloat64Summary{1.25, 9.75, 4.5}, g.PowerUsage)
	require.Equal(t, JobInt64Summary{11, 22, 17}, g.PCIeRxBandwidth)
	require.Equal(t, JobInt64Summary{41, 52, 47}, g.PCIeTxBandwidth)
	require.Equal(t, int64(31), g.PCIeReplays)
	require.Equal(t, int64(1000001), g.StartTime)
	require.Equal(t, int64(2000002), g.EndTime)
	require.Equal(t, JobInt32Summary{2, 95, 60}, g.SMUtilization)
	require.Equal(t, JobInt32Summary{3, 75, 40}, g.MemoryUtilization)
	require.Equal(t, JobInt32Summary{1000, 2000, 1500}, g.MemoryClock)
	require.Equal(t, JobInt32Summary{400, 1200, 800}, g.SMClock)
	require.Equal(t, uint(3), g.ECCSingleBit)
	require.Equal(t, uint(4), g.ECCDoubleBit)
	require.Equal(t, []int64{1234567}, g.XIDTimestamps)
	require.Equal(t, []JobProcessStats{{123, 3.5, 4.5}}, g.ComputeProcesses)
	require.Equal(t, []JobProcessStats{{456, 7.5, 8.5}}, g.GraphicsProcesses)
	require.Equal(t, int64(987654321), g.MaxGPUMemoryUsed)
	require.Equal(t, []int64{101, 102, 103, 104, 105, 106}, []int64{g.PowerViolationTime, g.ThermalViolationTime, g.ReliabilityViolationTime, g.BoardLimitViolationTime, g.LowUtilizationTime, g.SyncBoostTime})
	require.Len(t, g.HealthIncidents, 1)
	require.Equal(t, DCGM_HEALTH_RESULT_WARN, g.OverallHealth)
	require.Equal(t, []JobHealthIncident{{System: DCGM_HEALTH_WATCH_PCIE, Health: DCGM_HEALTH_RESULT_FAIL}}, g.HealthIncidents)
	require.Equal(t, DCGM_HEALTH_RESULT_FAIL, info.Summary.OverallHealth)
	require.Equal(t, []JobHealthIncident{{System: DCGM_HEALTH_WATCH_MEM, Health: DCGM_HEALTH_RESULT_PASS}}, info.Summary.HealthIncidents)
	info.GPUs[0].XIDTimestamps[0] = 0
	require.Equal(t, int64(1234667), info.Summary.XIDTimestamps[0])
	blank, err := jobBlankResponseForTest()
	require.NoError(t, err)
	require.True(t, IsInt64Blank(blank.Summary.EnergyConsumed))
	require.True(t, IsInt32Blank(int(blank.Summary.SMUtilization.Average)))
	require.True(t, IsFloat64Blank(blank.Summary.PowerUsage.Average))
}

func TestJobConversionBounds(t *testing.T) {
	for _, tc := range []struct {
		name                                  string
		gpu, xid, compute, graphics, incident int
	}{
		{"negative GPUs", -1, 0, 0, 0, 0},
		{"excess GPUs", int(MAX_NUM_DEVICES) + 1, 0, 0, 0, 0},
		{"negative XIDs", 0, -1, 0, 0, 0},
		{"excess XIDs", 0, 1000, 0, 0, 0},
		{"negative compute", 0, 0, -1, 0, 0},
		{"excess compute", 0, 0, 17, 0, 0},
		{"negative graphics", 0, 0, 0, -1, 0},
		{"excess graphics", 0, 0, 0, 17, 0},
		{"excess incidents", 0, 0, 0, 0, 1000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := jobResponseForTest(tc.gpu, tc.xid, tc.compute, tc.graphics, tc.incident)
			require.Error(t, err)
		})
	}
	info, err := jobResponseForTest(0, 0, 0, 0, 0)
	require.NoError(t, err)
	require.Empty(t, info.GPUs)
}

func TestJobConversionMaximums(t *testing.T) {
	info, err := jobMaxResponseForTest()
	require.NoError(t, err)
	require.Len(t, info.GPUs, int(MAX_NUM_DEVICES))
	require.Len(t, info.Summary.ComputeProcesses, 16)
	require.Len(t, info.Summary.GraphicsProcesses, 16)
	require.NotEmpty(t, info.Summary.XIDTimestamps)
	require.NotEmpty(t, info.Summary.HealthIncidents)
	_, err = jobInvalidGPUForTest()
	require.ErrorContains(t, err, "job GPU index 0")
}

func TestJobAPIsWithoutLibrary(t *testing.T) {
	previous := dcgmLibHandle
	dcgmLibHandle = nil
	t.Cleanup(func() { dcgmLibHandle = previous })
	cases := []struct {
		name, symbol string
		call         func() error
	}{
		{"watch", "dcgmWatchJobFields", func() error { return WatchJobFields(GroupHandle{}, time.Second, time.Minute, 0) }},
		{"start", "dcgmJobStartStats", func() error { return JobStart(GroupHandle{}, "valid-job") }},
		{"stop", "dcgmJobStopStats", func() error { return JobStop("valid-job") }},
		{"get", "dcgmJobGetStats", func() error { _, err := JobGetStats("valid-job"); return err }},
		{"remove", "dcgmJobRemove", func() error { return JobRemove("valid-job") }},
		{"remove all", "dcgmJobRemoveAll", JobRemoveAll},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			require.True(t, IsFeatureUnavailable(err))
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, err, &unavailable)
			require.Equal(t, tc.symbol, unavailable.Symbol)
		})
	}
}
