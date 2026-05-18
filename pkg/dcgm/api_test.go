/*
 * Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func setInitCounterForTest(t *testing.T, value int) {
	t.Helper()

	mux.Lock()
	previous := dcgmInitCounter
	dcgmInitCounter = value
	mux.Unlock()

	t.Cleanup(func() {
		mux.Lock()
		dcgmInitCounter = previous
		mux.Unlock()
	})
	require.Equal(t, 0, previous, "test requires clean DCGM init counter")
}

func initCounterForTest() int {
	mux.Lock()
	defer mux.Unlock()

	return dcgmInitCounter
}

func runWithDiscardedStderr(t *testing.T, f func()) {
	t.Helper()

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, devNull.Close())
	}()

	stderr := os.Stderr
	os.Stderr = devNull
	defer func() {
		os.Stderr = stderr
	}()

	f()
}

func TestInitWithNegativeCounterReturnsErrorWithoutCleanup(t *testing.T) {
	setInitCounterForTest(t, -2)

	cleanup, err := Init(Embedded)

	require.Error(t, err)
	require.Nil(t, cleanup)
	require.Contains(t, err.Error(), "shutdown() is called 2 times, before init()")
	require.Equal(t, -2, initCounterForTest())
}

func TestShutdownBeforeInitDoesNotCorruptCounter(t *testing.T) {
	setInitCounterForTest(t, 0)

	err := Shutdown()

	require.Error(t, err)
	require.Contains(t, err.Error(), "init() needs to be called before shutdown()")
	require.Equal(t, 0, initCounterForTest())

	err = Shutdown()

	require.Error(t, err)
	require.Contains(t, err.Error(), "init() needs to be called before shutdown()")
	require.Equal(t, 0, initCounterForTest())
}

func TestCleanupAfterManualShutdownDoesNotCorruptCounter(t *testing.T) {
	require.Equal(t, 0, initCounterForTest(), "test requires clean DCGM init counter")

	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	require.NotNil(t, cleanup)
	t.Cleanup(func() {
		mux.Lock()
		dcgmInitCounter = 0
		mux.Unlock()
	})

	require.NoError(t, Shutdown())
	require.Equal(t, 0, initCounterForTest())

	runWithDiscardedStderr(t, cleanup)
	require.Equal(t, 0, initCounterForTest())
}

func TestConcurrentShutdownBeforeInitDoesNotCorruptCounter(t *testing.T) {
	setInitCounterForTest(t, 0)

	const workers = 16
	start := make(chan struct{})
	errs := make(chan error, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs <- Shutdown()
		}()
	}

	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		require.Error(t, err)
		require.Contains(t, err.Error(), "init() needs to be called before shutdown()")
	}
	require.Equal(t, 0, initCounterForTest())
}

func TestGetEntityGroupEntities(t *testing.T) {
	withNvsdmMockConfig(t, "testdata/one_switch.yaml", func(t *testing.T) {
		teardownTest := setupTest(t)
		defer teardownTest(t)

		runOnlyWithLiveGPUs(t)

		// Get switch entities
		entities, err := GetEntityGroupEntities(FE_SWITCH)
		require.NoError(t, err)
		require.NotEmpty(t, entities)

		// Get nvlink entities
		nvlinkEntities, err := GetEntityGroupEntities(FE_LINK)
		require.NoError(t, err)
		require.NotEmpty(t, nvlinkEntities)
	})
}
