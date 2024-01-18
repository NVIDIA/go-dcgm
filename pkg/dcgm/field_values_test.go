/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetValuesSince(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	const gpu uint = 0

	// Create a group of fields
	const (
		xid int = iota
	)

	deviceFields := make([]Short, 1)
	deviceFields[xid] = DCGM_FI_DEV_XID_ERRORS

	fieldGroupName := fmt.Sprintf("fieldGroupName%d", rand.Uint64())
	fieldsGroup, err := FieldGroupCreate(fieldGroupName, deviceFields)
	assert.NoError(t, err)
	defer func() {
		_ = FieldGroupDestroy(fieldsGroup)
	}()

	t.Run("When there is no data return error", func(t *testing.T) {
		values, nextTime, err := GetValuesSince(GroupAllGPUs(),
			fieldsGroup, time.Time{})
		require.Error(t, err)
		assert.Empty(t, nextTime)
		assert.Len(t, values, 0)
	})

	t.Run("When there are a few entries", func(t *testing.T) {
		expectedNumberOfErrors := int64(43)
		expectedInjectedValuesCount := 0
		t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_XID_ERRORS", gpu)
		err = InjectFieldValue(gpu,
			DCGM_FI_DEV_XID_ERRORS,
			DCGM_FT_INT64,
			0,
			time.Now().Add(-time.Duration(5)*time.Second).UnixMicro(),
			expectedNumberOfErrors,
		)
		require.NoError(t, err)
		expectedInjectedValuesCount++
		for i := 4; i > 0; i-- {
			err = InjectFieldValue(gpu,
				DCGM_FI_DEV_XID_ERRORS,
				DCGM_FT_INT64,
				0,
				time.Now().Add(-time.Duration(i)*time.Second).UnixMicro(),
				int64(i),
			)
			require.NoError(t, err)
			expectedInjectedValuesCount++
		}
		// Force an update of the fields so that we can fetch initial values.
		err = UpdateAllFields()
		assert.NoError(t, err)
		values, nextTime, err := GetValuesSince(GroupAllGPUs(), fieldsGroup, time.Time{})
		assert.NoError(t, err)
		assert.Greater(t, nextTime, time.Time{})
		assert.Len(t, values, expectedInjectedValuesCount)
		assert.Equal(t, FE_GPU, values[0].EntityGroupId)
		assert.Equal(t, gpu, values[0].EntityId)
		assert.Equal(t, uint(DCGM_FI_DEV_XID_ERRORS), values[0].FieldId)
		assert.Equal(t, expectedNumberOfErrors, values[0].Int64())
		for i := 1; i < 5; i++ {
			assert.Equal(t, FE_GPU, values[i].EntityGroupId)
			assert.Equal(t, gpu, values[i].EntityId)
			assert.Equal(t, uint(DCGM_FI_DEV_XID_ERRORS), values[i].FieldId)
			assert.Equal(t, int64(5-i), values[i].Int64())
		}
	})
}
