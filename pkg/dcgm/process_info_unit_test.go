package dcgm

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTimeStringRunning verifies that a zero process timestamp is shown as a
// currently running process instead of as the Unix epoch.
func TestTimeStringRunning(t *testing.T) {
	if got := Time(0).String(); got != "Running" {
		t.Fatalf("Time(0).String() = %q, want Running", got)
	}
}

// TestTimeStringTimestamp verifies that a nonzero process timestamp uses Go's
// standard local-time representation.
func TestTimeStringTimestamp(t *testing.T) {
	timestamp := Time(1_700_000_000)
	assert.Equal(t, time.Unix(int64(timestamp), 0).String(), timestamp.String())
}

// TestProcessName verifies procfs name lookup for the current process and for a
// process ID that does not exist.
func TestProcessName(t *testing.T) {
	name, err := processName(uint(os.Getpid()))
	if err != nil || name == "" {
		t.Fatalf("processName(current PID) = %q, %v", name, err)
	}

	name, err = processName(math.MaxUint32)
	if err != nil || name != "" {
		t.Fatalf("processName(missing PID) = %q, %v, want empty nil", name, err)
	}
}

// TestProcessInfoFromResponse verifies that all per-GPU process metrics are
// converted correctly, including blank energy values and XID timestamps.
func TestProcessInfoFromResponse(t *testing.T) {
	const (
		pid  = 4242
		name = "training-process"
	)
	input := []processGPUInputForTest{
		{
			gpuID:                    3,
			energyConsumed:           12_345,
			startTime:                1_700_000_000_123_456,
			endTime:                  1_700_000_002_987_654,
			processSMUtilization:     61.4,
			processMemoryUtilization: 42.6,
			pcieRxBandwidth:          101,
			pcieTxBandwidth:          202,
			pcieReplays:              3,
			maxGPUMemoryUsed:         4_096,
			eccSingleBit:             4,
			eccDoubleBit:             5,
			smUtilization:            63,
			memoryUtilization:        47,
			memoryClock:              1_600,
			smClock:                  1_800,
			powerViolationTime:       11,
			thermalViolationTime:     12,
			reliabilityViolationTime: 13,
			boardLimitViolationTime:  14,
			lowUtilizationTime:       15,
			syncBoostTime:            16,
			reportedXIDCount:         2,
			xidTimestamps:            []int64{1_700_000_003, 1_700_000_004},
		},
		{
			gpuID:                    7,
			energyConsumed:           int64(dcgmInt64Blank),
			processSMUtilization:     75.5,
			processMemoryUtilization: 24.4,
		},
	}

	got, err := processInfoFromValuesForTest(pid, name, len(input), input)

	require.NoError(t, err)
	require.Len(t, got, len(input))

	firstGPU := got[0]
	assert.Equal(t, uint(3), firstGPU.GPU)
	assert.Equal(t, uint(pid), firstGPU.PID)
	assert.Equal(t, name, firstGPU.Name)
	assert.Equal(t, Time(1_700_000_000), firstGPU.ProcessUtilization.StartTime)
	assert.Equal(t, Time(1_700_000_002), firstGPU.ProcessUtilization.EndTime)
	require.NotNil(t, firstGPU.ProcessUtilization.EnergyConsumed)
	assert.Equal(t, uint64(12), *firstGPU.ProcessUtilization.EnergyConsumed)
	require.NotNil(t, firstGPU.ProcessUtilization.SmUtil)
	assert.Equal(t, float64(61), *firstGPU.ProcessUtilization.SmUtil)
	require.NotNil(t, firstGPU.ProcessUtilization.MemUtil)
	assert.Equal(t, float64(43), *firstGPU.ProcessUtilization.MemUtil)
	assert.Equal(t, PCIStatusInfo{Throughput: PCIThroughputInfo{
		Rx:      101,
		Tx:      202,
		Replays: 3,
	}}, firstGPU.PCI)
	assert.Equal(t, MemoryInfo{
		GlobalUsed: 4_096,
		ECCErrors:  ECCErrorsInfo{SingleBit: 4, DoubleBit: 5},
	}, firstGPU.Memory)
	assert.Equal(t, UtilizationInfo{GPU: 63, Memory: 47}, firstGPU.GpuUtilization)
	assert.Equal(t, ClockInfo{Cores: 1_800, Memory: 1_600}, firstGPU.Clocks)
	require.NotNil(t, firstGPU.Violations.Power)
	assert.Equal(t, uint64(11), *firstGPU.Violations.Power)
	require.NotNil(t, firstGPU.Violations.Thermal)
	assert.Equal(t, uint64(12), *firstGPU.Violations.Thermal)
	require.NotNil(t, firstGPU.Violations.Reliability)
	assert.Equal(t, uint64(13), *firstGPU.Violations.Reliability)
	require.NotNil(t, firstGPU.Violations.BoardLimit)
	assert.Equal(t, uint64(14), *firstGPU.Violations.BoardLimit)
	require.NotNil(t, firstGPU.Violations.LowUtilization)
	assert.Equal(t, uint64(15), *firstGPU.Violations.LowUtilization)
	require.NotNil(t, firstGPU.Violations.SyncBoost)
	assert.Equal(t, uint64(16), *firstGPU.Violations.SyncBoost)
	assert.Equal(t, XIDErrorInfo{
		NumErrors: 2,
		Timestamp: []uint64{1_700_000_003, 1_700_000_004},
	}, firstGPU.XIDErrors)

	secondGPU := got[1]
	assert.Equal(t, uint(7), secondGPU.GPU)
	assert.Equal(t, uint(pid), secondGPU.PID)
	assert.Equal(t, name, secondGPU.Name)
	assert.Equal(t, Time(0), secondGPU.ProcessUtilization.StartTime)
	assert.Equal(t, Time(0), secondGPU.ProcessUtilization.EndTime)
	require.NotNil(t, secondGPU.ProcessUtilization.EnergyConsumed)
	assert.Equal(t, uint64(0), *secondGPU.ProcessUtilization.EnergyConsumed)
	require.NotNil(t, secondGPU.ProcessUtilization.SmUtil)
	assert.Equal(t, float64(76), *secondGPU.ProcessUtilization.SmUtil)
	require.NotNil(t, secondGPU.ProcessUtilization.MemUtil)
	assert.Equal(t, float64(24), *secondGPU.ProcessUtilization.MemUtil)
	assert.Equal(t, PCIStatusInfo{}, secondGPU.PCI)
	assert.Equal(t, MemoryInfo{}, secondGPU.Memory)
	assert.Equal(t, UtilizationInfo{}, secondGPU.GpuUtilization)
	assert.Equal(t, ClockInfo{}, secondGPU.Clocks)
	require.NotNil(t, secondGPU.Violations.Power)
	assert.Equal(t, uint64(0), *secondGPU.Violations.Power)
	require.NotNil(t, secondGPU.Violations.Thermal)
	assert.Equal(t, uint64(0), *secondGPU.Violations.Thermal)
	require.NotNil(t, secondGPU.Violations.Reliability)
	assert.Equal(t, uint64(0), *secondGPU.Violations.Reliability)
	require.NotNil(t, secondGPU.Violations.BoardLimit)
	assert.Equal(t, uint64(0), *secondGPU.Violations.BoardLimit)
	require.NotNil(t, secondGPU.Violations.LowUtilization)
	assert.Equal(t, uint64(0), *secondGPU.Violations.LowUtilization)
	require.NotNil(t, secondGPU.Violations.SyncBoost)
	assert.Equal(t, uint64(0), *secondGPU.Violations.SyncBoost)
	assert.Equal(t, XIDErrorInfo{Timestamp: []uint64{}}, secondGPU.XIDErrors)
}

// TestProcessInfoFromResponseAcceptsLimits verifies that responses using every
// available native GPU and XID slot are accepted.
func TestProcessInfoFromResponseAcceptsLimits(t *testing.T) {
	maxGPUs, maxXIDs := processInfoLimitsForTest()
	input := make([]processGPUInputForTest, maxGPUs)
	for i := range input {
		input[i].gpuID = uint(i)
		input[i].reportedXIDCount = maxXIDs
		input[i].xidTimestamps = make([]int64, maxXIDs)
		for j := range input[i].xidTimestamps {
			input[i].xidTimestamps[j] = int64(i*maxXIDs + j)
		}
	}
	input[maxGPUs-1].xidTimestamps[maxXIDs-1] = 9_876

	got, err := processInfoFromValuesForTest(99, "limits", maxGPUs, input)

	require.NoError(t, err)
	require.Len(t, got, maxGPUs)
	assert.Equal(t, uint(maxGPUs-1), got[maxGPUs-1].GPU)
	require.Len(t, got[maxGPUs-1].XIDErrors.Timestamp, maxXIDs)
	assert.Equal(t, uint64(9_876), got[maxGPUs-1].XIDErrors.Timestamp[maxXIDs-1])
}

// TestProcessInfoFromResponseRejectsInvalidCounts verifies that reported GPU or
// XID counts outside the native arrays are rejected before they are indexed.
func TestProcessInfoFromResponseRejectsInvalidCounts(t *testing.T) {
	maxGPUs, maxXIDs := processInfoLimitsForTest()

	for _, gpuCount := range []int{-1, maxGPUs + 1} {
		t.Run(fmt.Sprintf("GPU count %d", gpuCount), func(t *testing.T) {
			got, err := processInfoFromValuesForTest(1, "invalid", gpuCount, nil)

			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf(
				"process info response contains %d GPUs, expected between 0 and %d",
				gpuCount,
				maxGPUs,
			))
		})
	}

	for _, xidCount := range []int{-1, maxXIDs + 1} {
		t.Run(fmt.Sprintf("XID count %d", xidCount), func(t *testing.T) {
			input := []processGPUInputForTest{{reportedXIDCount: xidCount}}

			got, err := processInfoFromValuesForTest(1, "invalid", len(input), input)

			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf(
				"process info response for GPU index 0 contains %d XID errors, expected between 0 and %d",
				xidCount,
				maxXIDs,
			))
		})
	}
}

// TestProcessInfoFromResponseEmpty verifies that a valid zero-GPU response
// becomes a non-nil, empty result for callers.
func TestProcessInfoFromResponseEmpty(t *testing.T) {
	got, err := processInfoFromValuesForTest(1, "empty", 0, nil)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}
