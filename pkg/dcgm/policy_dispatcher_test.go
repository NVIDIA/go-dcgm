//go:build linux && cgo

package dcgm

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUnknownConditionErrors(t *testing.T) {
	_, err := translateConditions(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least one policy condition")

	_, err = translateConditions([]PolicyCondition{PolicyCondition("unknown policy")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown policy condition")
}

func TestConcurrentPolicySubscribers(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1001)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	_, first, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	_, second, duplicateRegistration := dispatcher.addSubscription(group, xidCondition, 1)
	require.Nil(t, duplicateRegistration)

	violation := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 79},
	}
	dispatcher.deliver(registration.id, violation)

	assert.Equal(t, violation, receivePolicyViolation(t, first))
	assert.Equal(t, violation, receivePolicyViolation(t, second))
}

func TestPolicyDispatcherSlowConsumerDrop(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1002)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	_, slow, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)
	_, fast, _ := dispatcher.addSubscription(group, xidCondition, 1)

	first := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 31},
	}
	second := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 79},
	}

	dispatcher.deliver(registration.id, first)
	assert.Equal(t, first, receivePolicyViolation(t, fast))

	done := make(chan struct{})
	go func() {
		dispatcher.deliver(registration.id, second)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("deliver blocked on a full subscriber channel")
	}

	assert.Equal(t, second, receivePolicyViolation(t, fast))
	assert.Equal(t, uint64(1), dispatcher.dropped())
	assert.Equal(t, first, receivePolicyViolation(t, slow))
	assertNoPolicyViolation(t, slow)
}

func TestPolicyListenerLifecycle(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1003)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	firstID, first, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)
	_, second, _ := dispatcher.addSubscription(group, xidCondition, 1)

	closed, unregisters := dispatcher.removeSubscription(firstID)
	require.Equal(t, first, closed)
	assert.Empty(t, unregisters)
	close(closed)

	_, channelOpen := <-first
	assert.False(t, channelOpen)

	violation := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 43},
	}
	dispatcher.deliver(registration.id, violation)

	assert.Equal(t, violation, receivePolicyViolation(t, second))
}

func TestPolicyDispatcherDifferentGroupsIsolated(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	_, groupA, registrationA := dispatcher.addSubscription(
		policyTestGroupHandle(2001),
		xidCondition,
		1,
	)
	require.NotNil(t, registrationA)
	_, groupB, registrationB := dispatcher.addSubscription(
		policyTestGroupHandle(2002),
		xidCondition,
		1,
	)
	require.NotNil(t, registrationB)

	violation := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 95},
	}
	dispatcher.deliver(registrationA.id, violation)

	assert.Equal(t, violation, receivePolicyViolation(t, groupA))
	assertNoPolicyViolation(t, groupB)
}

func TestPolicyDispatcherFinalSubscriberReset(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1004)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	firstID, first, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	closed, unregisters := dispatcher.removeSubscription(firstID)
	require.Equal(t, first, closed)
	require.Len(t, unregisters, 1)
	assert.Equal(t, xidCondition, unregisters[0].condition)
	close(closed)

	dispatcher.mu.Lock()
	assert.Contains(t, dispatcher.registrations, registration.id)
	assert.Equal(t, xidCondition, dispatcher.registeredByGroup[group.GetHandle()])
	dispatcher.mu.Unlock()

	dispatcher.clearRegistrations(unregisters)

	_, second, nextRegistration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, nextRegistration)
	require.NotEqual(t, registration.id, nextRegistration.id)
	close(second)
}

func TestPolicyDispatcherClearRegistrationKeepsStateOnUnregisterFailure(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1005)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	firstID, first, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	closed, unregisters := dispatcher.removeSubscription(firstID)
	require.Equal(t, first, closed)
	require.Len(t, unregisters, 1)
	close(closed)

	_, next, nextRegistration := dispatcher.addSubscription(group, xidCondition, 1)
	require.Nil(t, nextRegistration)

	violation := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 48},
	}
	dispatcher.deliver(registration.id, violation)
	assert.Equal(t, violation, receivePolicyViolation(t, next))
}

func TestPolicyDispatcherRollbackSubscription(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1006)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	subID, ch, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	dispatcher.rollbackSubscription(subID, registration)

	_, channelOpen := <-ch
	assert.False(t, channelOpen)

	dispatcher.mu.Lock()
	assert.Empty(t, dispatcher.subscriptions)
	assert.Empty(t, dispatcher.registrations)
	assert.Empty(t, dispatcher.registeredByGroup)
	dispatcher.mu.Unlock()
}

func TestPolicyDispatcherDeliverUnknownCondition(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1007)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	_, ch, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	dispatcher.deliver(registration.id, PolicyViolation{Condition: PolicyCondition("unknown policy")})

	assertNoPolicyViolation(t, ch)
	assert.Equal(t, uint64(0), dispatcher.dropped())
}

func TestUnregisterErrorClearsLocalStateOnlyWhenRegistrationCannotBeLive(t *testing.T) {
	assert.True(t, unregisterErrorClearsLocalState(&Error{Code: -10}))
	assert.True(t, unregisterErrorClearsLocalState(&Error{Code: -21}))
	assert.False(t, unregisterErrorClearsLocalState(&Error{Code: -34}))
}

func TestPolicyReadNeedsDefaultSetupOnlyForMissingPolicyErrors(t *testing.T) {
	assert.True(t, policyReadNeedsDefaultSetup(&Error{Code: -31}))
	assert.True(t, policyReadNeedsDefaultSetup(&Error{Code: -5}))

	assert.False(t, policyReadNeedsDefaultSetup(&Error{Code: -8}))
	assert.False(t, policyReadNeedsDefaultSetup(&Error{Code: -21}))
	assert.False(t, policyReadNeedsDefaultSetup(assert.AnError))
}

func TestPolicyDispatcherDeliverDuringUnsubscribeRace(t *testing.T) {
	dispatcher := newPolicyDispatcher()
	group := policyTestGroupHandle(1008)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	subID, ch, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	violation := PolicyViolation{
		Condition: XidPolicy,
		Data:      XidPolicyCondition{ErrNum: 79},
	}

	start := make(chan struct{})
	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for {
				select {
				case <-stop:
					return
				default:
					dispatcher.deliver(registration.id, violation)
				}
			}
		}()
	}

	close(start)
	time.Sleep(10 * time.Millisecond)

	closed, unregisters := dispatcher.removeSubscription(subID)
	require.Equal(t, ch, closed)
	require.Len(t, unregisters, 1)
	close(closed)

	close(stop)
	wg.Wait()
}

func TestPolicyConfigsForListenPreserveExistingThresholds(t *testing.T) {
	status := &PolicyStatus{
		Action:     PolicyActionGPUReset,
		Validation: PolicyValidationShort,
		Conditions: map[PolicyCondition]interface{}{
			ThermalPolicy: uint32(85),
		},
	}

	configs, needsUpdate := policyConfigsForListen(status, []PolicyCondition{ThermalPolicy, XidPolicy})

	require.True(t, needsUpdate)
	require.Len(t, configs, 2)
	assert.Equal(t, ThermalPolicy, configs[0].Condition)
	require.NotNil(t, configs[0].MaxTemperature)
	assert.Equal(t, uint32(85), *configs[0].MaxTemperature)
	require.NotNil(t, configs[0].Action)
	assert.Equal(t, PolicyActionGPUReset, *configs[0].Action)
	require.NotNil(t, configs[0].Validation)
	assert.Equal(t, PolicyValidationShort, *configs[0].Validation)
	assert.Equal(t, XidPolicy, configs[1].Condition)
}

func TestPolicyConfigsForListenNoopWhenRequestedConditionsExist(t *testing.T) {
	status := &PolicyStatus{
		Conditions: map[PolicyCondition]interface{}{
			XidPolicy: true,
		},
	}

	configs, needsUpdate := policyConfigsForListen(status, []PolicyCondition{XidPolicy})

	assert.False(t, needsUpdate)
	assert.Empty(t, configs)
}

func TestPolicyViolationDropCount(t *testing.T) {
	previousCallbacks := policyCallbacks
	dispatcher := newPolicyDispatcher()
	policyCallbacks = dispatcher
	t.Cleanup(func() {
		policyCallbacks = previousCallbacks
	})

	group := policyTestGroupHandle(1009)
	xidCondition, ok := policyConditionMask(XidPolicy)
	require.True(t, ok)

	_, ch, registration := dispatcher.addSubscription(group, xidCondition, 1)
	require.NotNil(t, registration)

	dispatcher.deliver(registration.id, PolicyViolation{Condition: XidPolicy})
	dispatcher.deliver(registration.id, PolicyViolation{Condition: XidPolicy})

	assert.Equal(t, uint64(1), PolicyViolationDropCount())
	_ = receivePolicyViolation(t, ch)
}

func policyTestGroupHandle(id uintptr) GroupHandle {
	var group GroupHandle
	group.SetHandle(id)
	return group
}

func receivePolicyViolation(t *testing.T, ch <-chan PolicyViolation) PolicyViolation {
	t.Helper()

	select {
	case violation, ok := <-ch:
		require.True(t, ok, "policy violation channel closed")
		return violation
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for policy violation")
		return PolicyViolation{}
	}
}

func assertNoPolicyViolation(t *testing.T, ch <-chan PolicyViolation) {
	t.Helper()

	select {
	case violation, ok := <-ch:
		if ok {
			t.Fatalf("unexpected policy violation: %+v", violation)
		}
	case <-time.After(50 * time.Millisecond):
	}
}
