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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigDeviceProfileNamesStandalone(t *testing.T) {
	// Setup test environment
	teardown := setupTest(t)
	defer teardown(t)

	// Create one fake GPU
	gpuIDs, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)
	require.Len(t, gpuIDs, 1, "Expected 1 fake GPU to be created")

	// Create one GPU instance on the fake GPU
	gpuInstanceMap, err := withInjectionGPUInstances(t, gpuIDs[0], 1)
	require.NoError(t, err)
	require.Len(t, gpuInstanceMap, 1, "Expected 1 fake GPU instance to be created")

	// Get the GPU instance IDs
	gpuInstanceIDs := make([]uint, 0, len(gpuInstanceMap))
	for instanceID := range gpuInstanceMap {
		gpuInstanceIDs = append(gpuInstanceIDs, instanceID)
	}

	// Create one compute instance per GPU instance
	ciToGiMap, err := withInjectionComputeInstances(t, gpuInstanceIDs, len(gpuInstanceIDs))
	require.NoError(t, err)
	require.Len(t, ciToGiMap, len(gpuInstanceIDs), "Expected one compute instance per GPU instance")

	// Get the compute instance IDs
	computeInstanceIds := make([]uint, 0, len(ciToGiMap))
	for ciId := range ciToGiMap {
		computeInstanceIds = append(computeInstanceIds, ciId)
	}

	// Verify profile names for both GPU instances and compute instances
	verifyProfileNames(t, gpuInstanceIDs, true)      // verify GPU instances
	verifyProfileNames(t, computeInstanceIds, false) // verify compute instances
}

// verifyProfileNames verifies that the MIG profile names exist for the given entities
func verifyProfileNames(tb testing.TB, entityIds []uint, isGpuInstance bool) {
	tb.Helper()

	// Create entity list for the query
	entities := make([]GroupEntityPair, 0, len(entityIds))
	for _, entityId := range entityIds {
		entity := GroupEntityPair{
			EntityId: entityId,
		}
		if isGpuInstance {
			entity.EntityGroupId = FE_GPU_I
		} else {
			entity.EntityGroupId = FE_GPU_CI
		}
		entities = append(entities, entity)
	}

	// Get the latest values for DCGM_FI_DEV_NAME field
	values, err := EntitiesGetLatestValues(entities, []Short{DCGM_FI_DEV_NAME}, DCGM_FV_FLAG_LIVE_DATA)
	require.NoError(tb, err)

	// Define expected profile names
	expectedFakeName := "1fc.1g.4gb"
	if isGpuInstance {
		expectedFakeName = "1fg.4gb"
	}

	// Verify each entity has the correct profile name
	for i := range values {
		assert.Equal(tb, expectedFakeName, values[i].String(),
			"Fake profile name appears to be wrong for entity %d. Expected '%s', found '%s'",
			values[i].EntityID, expectedFakeName, values[i].String())
	}
}
