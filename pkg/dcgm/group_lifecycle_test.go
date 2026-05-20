//go:build linux && cgo

package dcgm

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func requireNoGroupCapErrorForTest(t *testing.T, iter int, err error) {
	t.Helper()
	require.Error(t, err)

	require.Falsef(t,
		isGroupCapErrorForTest(err),
		"iter %d: group leak detected; CreateGroup hit DCGM_MAX_NUM_GROUPS: %v", iter, err)
}

func isGroupCapErrorForTest(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "max limit") || strings.Contains(msg, "max number of groups") ||
		strings.Contains(msg, "max group")
}

func skipIfPidWatchRequiresRoot(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "non-root") {
		t.Skipf("PID field watches require root on this DCGM host: %v", err)
	}
}

func TestGroupCapErrorMatcherRecognizesDCGMMaxLimit(t *testing.T) {
	require.True(t, isGroupCapErrorForTest(errors.New("Max limit reached for the object")))
	require.True(t, isGroupCapErrorForTest(errors.New("Add Group: Max number of groups already configured")))
	require.False(t, isGroupCapErrorForTest(errors.New("This GPU is not supported by DCGM")))
}

func TestHealthCheckByGpuIdDoesNotLeakGroupOnError(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	const invalidGPU uint = 1 << 30
	for i := 0; i < 70; i++ {
		_, err := healthCheckByGpuId(invalidGPU)
		requireNoGroupCapErrorForTest(t, i, err)
	}
}

func TestWatchFieldsDoesNotLeakGroupOnError(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	fieldGroup, err := FieldGroupCreate("watch-fields-leak", []Short{DCGM_FI_DEV_GPU_TEMP})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, FieldGroupDestroy(fieldGroup))
	}()

	const invalidGPU uint = 1 << 30
	for i := 0; i < 70; i++ {
		_, err := WatchFields(invalidGPU, fieldGroup, "watch-fields-leak")
		requireNoGroupCapErrorForTest(t, i, err)
	}
}

func TestWatchPidFieldsDoesNotLeakGroupOnError(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	const invalidGPU uint = 1 << 30
	for i := 0; i < 70; i++ {
		_, err := watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples, invalidGPU)
		requireNoGroupCapErrorForTest(t, i, err)
	}
}

func TestWatchPidFieldsDoesNotLeakGroupOnNonRootWatchError(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	// This regression is intentionally environment-specific: it covers
	// dcgmWatchPidFields returning DCGM_ST_REQUIRES_ROOT on non-root hosts.
	// Root hosts cover PID-watch success ownership in
	// TestWatchPidFieldsDoesNotDestroyGroupOnSuccess.
	for i := 0; i < 70; i++ {
		group, err := watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples, gpus[0])
		if err == nil {
			require.NoError(t, DestroyGroup(group))
			t.Skip("PID field watches succeeded; non-root watch-error cleanup path is not reproducible on this host")
		}

		requireNoGroupCapErrorForTest(t, i, err)

		var dcgmErr *Error
		require.ErrorAs(t, err, &dcgmErr)
	}
}

func TestWatchFieldsDoesNotDestroyGroupOnSuccess(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	fieldGroup, err := FieldGroupCreate("watch-fields-success", []Short{DCGM_FI_DEV_GPU_TEMP})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, FieldGroupDestroy(fieldGroup))
	}()

	group, err := WatchFields(gpus[0], fieldGroup, "watch-fields-success")
	require.NoError(t, err)
	require.NoError(t, DestroyGroup(group))
}

func TestWatchPidFieldsDoesNotDestroyGroupOnSuccess(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	group, err := watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples, gpus[0])
	skipIfPidWatchRequiresRoot(t, err)
	require.NoError(t, err)
	require.NoError(t, DestroyGroup(group))
}

func TestWatchFieldsReturnsUpdateAllFieldsErrorAndDestroysGroup(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	fieldGroup, err := FieldGroupCreate("watch-fields-update-fail", []Short{DCGM_FI_DEV_XID_ERRORS})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, FieldGroupDestroy(fieldGroup))
	}()

	updateErr := errors.New("forced update error")
	oldUpdateAllFields := updateAllFields
	updateAllFields = func() error { return updateErr }
	t.Cleanup(func() { updateAllFields = oldUpdateAllFields })

	for i := 0; i < 70; i++ {
		_, err = WatchFields(gpus[0], fieldGroup, "watch-fields-update-fail")
		requireNoGroupCapErrorForTest(t, i, err)
		require.ErrorIs(t, err, updateErr)
	}
}

func TestWatchPidFieldsReturnsUpdateAllFieldsErrorAndDestroysGroup(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	updateErr := errors.New("forced update error")
	oldUpdateAllFields := updateAllFields
	updateAllFields = func() error { return updateErr }
	t.Cleanup(func() { updateAllFields = oldUpdateAllFields })

	for i := 0; i < 70; i++ {
		_, err = watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples, gpus[0])
		skipIfPidWatchRequiresRoot(t, err)
		requireNoGroupCapErrorForTest(t, i, err)
		require.ErrorIs(t, err, updateErr)
	}
}
