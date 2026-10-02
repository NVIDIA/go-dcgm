/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestValidateFieldPolicySpec(t *testing.T) {
	valid := FieldPolicySpec{
		FieldID:  DCGM_FI_DEV_GPU_TEMP,
		Name:     "temperature",
		Operator: FieldPolicyGreaterThan,
	}
	require.NoError(t, validateFieldPolicySpec(valid))

	tests := []struct {
		name string
		spec FieldPolicySpec
	}{
		{name: "zero field", spec: FieldPolicySpec{Name: "policy", Operator: FieldPolicyGreaterThan}},
		{name: "empty name", spec: FieldPolicySpec{FieldID: DCGM_FI_DEV_GPU_TEMP, Operator: FieldPolicyGreaterThan}},
		{name: "NUL name", spec: FieldPolicySpec{FieldID: DCGM_FI_DEV_GPU_TEMP, Name: "bad\x00name", Operator: FieldPolicyGreaterThan}},
		{name: "unknown operator", spec: FieldPolicySpec{FieldID: DCGM_FI_DEV_GPU_TEMP, Name: "policy", Operator: FieldPolicyOperator(99)}},
		{name: "too many entities", spec: FieldPolicySpec{
			FieldID:  DCGM_FI_DEV_GPU_TEMP,
			Name:     "policy",
			Operator: FieldPolicyGreaterThan,
			Entities: make([]GroupEntityPair, maxFieldPolicyEntitiesForTest()+1),
		}},
	}
	if uint64(^uint(0)) > uint64(^uint32(0)) {
		rateLimit := uint(^uint32(0))
		rateLimit++
		tests = append(tests, struct {
			name string
			spec FieldPolicySpec
		}{name: "rate limit exceeds C uint", spec: FieldPolicySpec{
			FieldID:          DCGM_FI_DEV_GPU_TEMP,
			Name:             "policy",
			Operator:         FieldPolicyGreaterThan,
			RateLimitSeconds: rateLimit,
		}})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Error(t, validateFieldPolicySpec(test.spec))
		})
	}
}

func TestFieldPolicyReadbackMapping(t *testing.T) {
	policy := createTestFieldPolicyReadback()
	assert.Equal(t, uint64(41), policy.PolicyID)
	assert.Equal(t, DCGM_FI_DEV_GPU_TEMP, policy.FieldID)
	assert.Equal(t, "temperature", policy.Name)
	assert.Equal(t, 81.5, policy.CurrentValue)
	assert.Equal(t, 80.0, policy.Threshold)
	assert.Equal(t, FieldPolicyGreaterThan, policy.Operator)
	assert.Equal(t, []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: 3}}, policy.Entities)
	assert.Equal(t, FieldPolicyChannelCallback, policy.Channels)
	assert.True(t, policy.Enabled)
	assert.Equal(t, uint(7), policy.RateLimitSeconds)
	assert.Equal(t, time.UnixMicro(1_700_000_000_123_456), policy.CreatedAt)
	assert.Equal(t, uint(9), policy.ViolationCount)
	assert.Equal(t, time.UnixMicro(1_700_000_001_123_456), policy.LastTriggeredTimestamp)
}

func TestFieldPolicyV47LayoutSemantics(t *testing.T) {
	specVersion, expectedSpecVersion, infoVersion, expectedInfoVersion := fieldPolicyLayoutForTest()
	assert.Equal(t, expectedSpecVersion, specVersion)
	assert.Equal(t, expectedInfoVersion, infoVersion)
	assert.Equal(t, uint(4_294_967_295), fieldPolicyRateLimitForTest())
}

func TestFieldPolicyUnavailableSymbol(t *testing.T) {
	err := requireFieldPolicySymbol("dcgmPolicyUnavailableForTest")
	require.Error(t, err)
	assert.True(t, IsFeatureUnavailable(err))
}

func TestFieldPolicyUnavailableResult(t *testing.T) {
	err := fieldPolicyFunctionNotFoundErrorForTest()
	require.Error(t, err)
	assert.True(t, IsFeatureUnavailable(err))

	var unavailable *FeatureUnavailableError
	require.ErrorAs(t, err, &unavailable)
	assert.Equal(t, fieldPolicyFeature, unavailable.Feature)
	assert.Equal(t, "dcgmPolicyGetAll", unavailable.Symbol)

	var dcgmErr *Error
	require.ErrorAs(t, err, &dcgmErr)
	assert.Equal(t, fieldPolicyFunctionNotFoundCodeForTest(), int(dcgmErr.Code))
}

func TestFieldPolicyCountLimit(t *testing.T) {
	require.NoError(t, validateFieldPolicyCount(MaxFieldPolicies))
	assert.Error(t, validateFieldPolicyCount(MaxFieldPolicies+1))
}

// TestGetAllFieldPoliciesWithCalls checks list retrieval, capacity changes, and native read errors.
func TestGetAllFieldPoliciesWithCalls(t *testing.T) {
	callErr := errors.New("DCGM failed")
	policy := FieldPolicy{PolicyID: 7}
	for _, tt := range []struct {
		name      string
		count     int
		countErr  error
		list      []FieldPolicy
		reported  int
		listErr   error
		wantError string
	}{
		{name: "count error", countErr: callErr, wantError: "DCGM failed"},
		{name: "empty"},
		{name: "count exceeds limit", count: MaxFieldPolicies + 1, wantError: "maximum supported"},
		{name: "list error", count: 1, listErr: callErr, wantError: "DCGM failed"},
		{name: "list grows", count: 1, reported: 2, wantError: "list grew"},
		{name: "success", count: 1, reported: 1, list: []FieldPolicy{policy}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockFieldPolicyListCalls(gomock.NewController(t))
			mock.EXPECT().fieldPolicyCount().Return(tt.count, tt.countErr)
			if tt.count > 0 && tt.count <= MaxFieldPolicies && tt.countErr == nil {
				mock.EXPECT().fieldPolicyList(tt.count).Return(tt.list, tt.reported, tt.listErr)
			}
			got, err := getAllFieldPoliciesWithCalls(mock)
			if tt.wantError != "" {
				require.ErrorContains(t, err, tt.wantError)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.list, got)
		})
	}
}

// TestSubscribeFieldPolicyWithCallsRollsBack checks local cleanup after native registration fails.
func TestSubscribeFieldPolicyWithCallsRollsBack(t *testing.T) {
	registrationErr := errors.New("registration failed")
	for _, tt := range []struct {
		name      string
		register  bool
		wantError string
	}{
		{name: "registration failure", register: true, wantError: "registration failed"},
		{name: "duplicate subscription", wantError: "already registered"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			registry := newFieldPolicySubscriptionRegistry()
			if !tt.register {
				_, err := registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 1)
				require.NoError(t, err)
			}
			mock := NewMockFieldPolicyEventCalls(gomock.NewController(t))
			if tt.register {
				mock.EXPECT().registerFieldPolicyNative(Short(203), uint64(4), uint64(1)).Return(registrationErr)
			}
			_, err := subscribeFieldPolicyWithCalls(mock, registry, 203, 4, 1)
			require.ErrorContains(t, err, tt.wantError)
			if tt.register {
				require.Empty(t, registry.subscriptions)
				require.Empty(t, registry.byKey)
			}
		})
	}
}

// TestSubscribeFieldPolicyWithCallsClosesThroughAPI checks that Close unregisters
// the DCGM callback, removes the local subscription, and closes its channel.
func TestSubscribeFieldPolicyWithCallsClosesThroughAPI(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	m := NewMockFieldPolicyEventCalls(gomock.NewController(t))
	m.EXPECT().registerFieldPolicyNative(Short(203), uint64(4), uint64(1)).Return(nil)
	m.EXPECT().unregisterFieldPolicyNative(Short(203), uint64(4)).Return(nil)

	subscription, err := subscribeFieldPolicyWithCalls(m, registry, 203, 4, 1)
	require.NoError(t, err)
	require.NoError(t, subscription.Close())
	require.Empty(t, registry.subscriptions)
	require.Empty(t, registry.byKey)
	_, open := <-subscription.Violations
	require.False(t, open)
}

func TestCopyGoStringToCCharsPreservesUTF8(t *testing.T) {
	assert.Equal(t, "温度-policy", fieldPolicyNameForTest("温度-policy"))
}

func TestMicrosecondsOrZero(t *testing.T) {
	assert.Zero(t, microsecondsOrZero(time.Time{}))
	assert.Equal(t, int64(1_700_000_000_123_456), microsecondsOrZero(time.UnixMicro(1_700_000_000_123_456)))
}

func TestFieldPolicySubscriptionClose(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	subscription, err := registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 1)
	require.NoError(t, err)

	violation := FieldPolicyViolation{PolicyID: 4, FieldID: 203}
	registry.deliver(subscription.id, violation)
	require.Equal(t, violation, <-subscription.Violations)

	unregisterCalls := 0
	subscription.unregister = func() error {
		unregisterCalls++
		return nil
	}
	require.NoError(t, subscription.Close())
	require.Equal(t, 1, unregisterCalls)
	require.NoError(t, subscription.Close())
	require.Equal(t, 1, unregisterCalls)

	_, open := <-subscription.Violations
	assert.False(t, open)
	registry.deliver(subscription.id, violation)
}

func TestFieldPolicySubscriptionCloseClearsLocalStateAfterDCGMDisconnect(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	key := fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}
	subscription, err := registry.add(key, 1)
	require.NoError(t, err)

	subscription.unregister = fieldPolicyUninitializedErrorForTest
	require.NoError(t, subscription.Close())
	_, err = registry.add(key, 1)
	require.NoError(t, err)
}

func TestFieldPolicySubscriptionRejectsDuplicateKey(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	_, err := registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 1)
	require.NoError(t, err)
	_, err = registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 1)
	assert.ErrorContains(t, err, "already registered")
}

func TestFieldPolicySubscriptionDeliveryDoesNotBlock(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	subscription, err := registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 1)
	require.NoError(t, err)
	registry.deliver(subscription.id, FieldPolicyViolation{PolicyID: 4})

	done := make(chan struct{})
	go func() {
		registry.deliver(subscription.id, FieldPolicyViolation{PolicyID: 4})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("field-policy delivery blocked on a full subscriber channel")
	}
}

func TestFieldPolicySubscriptionConcurrentCloseAndDelivery(t *testing.T) {
	registry := newFieldPolicySubscriptionRegistry()
	subscription, err := registry.add(fieldPolicySubscriptionKey{fieldID: 203, policyID: 4}, 8)
	require.NoError(t, err)

	var workers sync.WaitGroup
	for i := 0; i < 8; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for j := 0; j < 100; j++ {
				registry.deliver(subscription.id, FieldPolicyViolation{PolicyID: 4})
			}
		}()
	}
	subscription.unregister = func() error { return nil }
	require.NoError(t, subscription.Close())
	workers.Wait()
	for range subscription.Violations {
	}
	_, open := <-subscription.Violations
	assert.False(t, open)
}
