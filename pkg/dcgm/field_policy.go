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

/*
#include <stdint.h>
#include "dcgm_agent.h"
#include "dcgm_structs.h"

extern int fieldPolicyViolationNotify(dcgm_policy_violation_t *violation, uint64_t userData);
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	// MaxFieldPolicies is the largest number of field policies DCGM v4.7 can return.
	MaxFieldPolicies = int(C.DCGM_POLICY_MAX_POLICIES)

	// MaxFieldPolicyViolations is the largest number of policy violations DCGM v4.7 returns at once.
	MaxFieldPolicyViolations = int(C.DCGM_POLICY_MAX_VIOLATIONS)
)

// FieldPolicyOperator compares the current field value with a policy threshold.
type FieldPolicyOperator uint32

const (
	// FieldPolicyGreaterThan triggers when the field value is greater than the threshold.
	FieldPolicyGreaterThan FieldPolicyOperator = C.DCGM_POLICY_OP_GT
	// FieldPolicyGreaterThanOrEqual triggers when the field value is at least the threshold.
	FieldPolicyGreaterThanOrEqual FieldPolicyOperator = C.DCGM_POLICY_OP_GE
	// FieldPolicyLessThan triggers when the field value is less than the threshold.
	FieldPolicyLessThan FieldPolicyOperator = C.DCGM_POLICY_OP_LT
	// FieldPolicyLessThanOrEqual triggers when the field value is at most the threshold.
	FieldPolicyLessThanOrEqual FieldPolicyOperator = C.DCGM_POLICY_OP_LE
	// FieldPolicyEqual triggers when the field value equals the threshold.
	FieldPolicyEqual FieldPolicyOperator = C.DCGM_POLICY_OP_EQ
	// FieldPolicyNotEqual triggers when the field value differs from the threshold.
	FieldPolicyNotEqual FieldPolicyOperator = C.DCGM_POLICY_OP_NE
	// FieldPolicyChanged triggers when the field value changes.
	FieldPolicyChanged FieldPolicyOperator = C.DCGM_POLICY_OP_CHANGED
)

// FieldPolicyChannels selects where DCGM delivers a field-policy violation.
type FieldPolicyChannels uint32

const (
	// FieldPolicyChannelConsole writes a violation to the host engine console.
	FieldPolicyChannelConsole FieldPolicyChannels = C.DCGM_POLICY_CHANNEL_CONSOLE
	// FieldPolicyChannelFile writes a violation to the host engine policy log.
	FieldPolicyChannelFile FieldPolicyChannels = C.DCGM_POLICY_CHANNEL_FILE
	// FieldPolicyChannelCallback delivers a violation to registered callbacks.
	FieldPolicyChannelCallback FieldPolicyChannels = C.DCGM_POLICY_CHANNEL_CALLBACK
)

// FieldPolicySpec describes a field-policy to create.
type FieldPolicySpec struct {
	FieldID          Short
	Name             string
	Threshold        float64
	Operator         FieldPolicyOperator
	Entities         []GroupEntityPair
	Channels         FieldPolicyChannels
	Enabled          bool
	RateLimitSeconds uint
}

// FieldPolicy is the current state of a field-policy returned by DCGM.
type FieldPolicy struct {
	PolicyID               uint64
	FieldID                Short
	Name                   string
	CurrentValue           float64
	Threshold              float64
	Operator               FieldPolicyOperator
	Entities               []GroupEntityPair
	Channels               FieldPolicyChannels
	Enabled                bool
	RateLimitSeconds       uint
	CreatedAt              time.Time
	ViolationCount         uint
	LastTriggeredTimestamp time.Time
}

// FieldPolicyViolation records one field-policy threshold violation.
type FieldPolicyViolation struct {
	PolicyID     uint64
	PolicyName   string
	FieldName    string
	FieldID      Short
	Entity       GroupEntityPair
	CurrentValue float64
	Threshold    float64
	Operator     FieldPolicyOperator
	Channels     FieldPolicyChannels
	Timestamp    time.Time
}

// FeatureUnavailableError reports an API that is not available in the loaded DCGM library or host engine.
type FeatureUnavailableError struct {
	Feature string
	Symbol  string
	Cause   error
}

func (e *FeatureUnavailableError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s is unavailable (%s): %v", e.Feature, e.Symbol, e.Cause)
	}
	return fmt.Sprintf("%s is unavailable (%s)", e.Feature, e.Symbol)
}

// Unwrap returns the DCGM error reported by the host engine, when it is available.
func (e *FeatureUnavailableError) Unwrap() error {
	return e.Cause
}

// IsFeatureUnavailable reports whether err identifies an unavailable DCGM feature.
func IsFeatureUnavailable(err error) bool {
	var unavailable *FeatureUnavailableError
	return errors.As(err, &unavailable)
}

const fieldPolicyFeature = "DCGM field-policy API"

func requireFieldPolicySymbol(symbol string) error {
	if dcgmSymbolAvailable(symbol) {
		return nil
	}
	return &FeatureUnavailableError{Feature: fieldPolicyFeature, Symbol: symbol}
}

func fieldPolicyResult(symbol string, result C.dcgmReturn_t) error {
	return dcgmFeatureResult(fieldPolicyFeature, symbol, result, C.GoString(C.errorString(result)))
}

func dcgmFeatureResult(feature, symbol string, result C.dcgmReturn_t, message string) error {
	if result == C.DCGM_ST_OK {
		return nil
	}
	dcgmErr := &Error{msg: message, Code: result}
	if result == C.DCGM_ST_FUNCTION_NOT_FOUND {
		return &FeatureUnavailableError{Feature: feature, Symbol: symbol, Cause: dcgmErr}
	}
	return dcgmErr
}

func validateFieldPolicySpec(spec FieldPolicySpec) error {
	if spec.FieldID == 0 {
		return errors.New("field policy field ID must not be zero")
	}
	if err := validateCString("field policy name", spec.Name, int(C.DCGM_MAX_STR_LENGTH)); err != nil {
		return err
	}
	if spec.Operator > FieldPolicyChanged {
		return fmt.Errorf("unsupported field policy operator: %d", spec.Operator)
	}
	if len(spec.Entities) > int(C.DCGM_MAX_NUM_DEVICES) {
		return fmt.Errorf("field policy has %d entities; DCGM supports at most %d", len(spec.Entities), C.DCGM_MAX_NUM_DEVICES)
	}
	if spec.RateLimitSeconds > math.MaxUint32 {
		return fmt.Errorf("field policy rate limit %d exceeds the DCGM maximum of %d seconds", spec.RateLimitSeconds, uint(math.MaxUint32))
	}
	return nil
}

func validateCString(name, value string, maxSize int) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty", name)
	}
	if strings.IndexByte(value, 0) >= 0 {
		return fmt.Errorf("%s must not contain a NUL byte", name)
	}
	if len(value) >= maxSize {
		return fmt.Errorf("%s is too long; maximum length is %d bytes", name, maxSize-1)
	}
	return nil
}

func cFieldPolicyEntities(entities []GroupEntityPair) []C.dcgmGroupEntityPair_t {
	if len(entities) == 0 {
		return nil
	}

	cEntities := make([]C.dcgmGroupEntityPair_t, len(entities))
	for i, entity := range entities {
		cEntities[i].entityGroupId = C.dcgm_field_entity_group_t(entity.EntityGroupId)
		cEntities[i].entityId = C.uint(entity.EntityId)
	}
	return cEntities
}

// CreateFieldPolicy creates a field-policy and returns its DCGM-assigned identifier.
func CreateFieldPolicy(spec FieldPolicySpec) (uint64, error) {
	const symbol = "dcgmPolicyCreate"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return 0, err
	}
	if err := validateFieldPolicySpec(spec); err != nil {
		return 0, err
	}

	cName := C.CString(spec.Name)
	defer freeCString(cName)
	cEntities := cFieldPolicyEntities(spec.Entities)
	var entities *C.dcgmGroupEntityPair_t
	if len(cEntities) > 0 {
		entities = &cEntities[0]
	}
	var policyID C.uint64_t
	result := C.dcgmPolicyCreate(
		handle.handle,
		C.ushort(spec.FieldID),
		cName,
		C.double(spec.Threshold),
		C.dcgmPolicyOperator_t(spec.Operator),
		entities,
		C.uint(len(cEntities)),
		C.uint(spec.Channels),
		C.uint(boolToUint(spec.Enabled)),
		C.uint(spec.RateLimitSeconds),
		&policyID,
	)
	if err := fieldPolicyResult(symbol, result); err != nil {
		return 0, err
	}
	return uint64(policyID), nil
}

// ModifyFieldPolicy replaces the editable settings of an existing field-policy.
func ModifyFieldPolicy(policy FieldPolicy) error {
	const symbol = "dcgmPolicyModify"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	if policy.PolicyID == 0 {
		return errors.New("field policy ID must not be zero")
	}
	cPolicy, err := cFieldPolicyInfo(policy)
	if err != nil {
		return err
	}
	return fieldPolicyResult(symbol, C.dcgmPolicyModify(handle.handle, &cPolicy))
}

// DeleteFieldPolicy deletes one field-policy.
func DeleteFieldPolicy(fieldID Short, policyID uint64) error {
	const symbol = "dcgmPolicyDelete"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	return fieldPolicyResult(symbol, C.dcgmPolicyDelete(handle.handle, C.ushort(fieldID), C.uint64_t(policyID), 0))
}

// DeleteAllFieldPolicies deletes every field-policy registered with DCGM.
func DeleteAllFieldPolicies() error {
	const symbol = "dcgmPolicyDelete"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	return fieldPolicyResult(symbol, C.dcgmPolicyDelete(handle.handle, 0, 0, 1))
}

// EnableFieldPolicy enables one field-policy.
func EnableFieldPolicy(fieldID Short, policyID uint64) error {
	const symbol = "dcgmPolicyEnable"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	return fieldPolicyResult(symbol, C.dcgmPolicyEnable(handle.handle, C.ushort(fieldID), C.uint64_t(policyID)))
}

// DisableFieldPolicy disables one field-policy.
func DisableFieldPolicy(fieldID Short, policyID uint64) error {
	const symbol = "dcgmPolicyDisable"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	return fieldPolicyResult(symbol, C.dcgmPolicyDisable(handle.handle, C.ushort(fieldID), C.uint64_t(policyID)))
}

// GetFieldPolicy returns one field-policy.
func GetFieldPolicy(fieldID Short, policyID uint64) (FieldPolicy, error) {
	const symbol = "dcgmPolicyGetOne"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return FieldPolicy{}, err
	}

	var cPolicy C.dcgmPolicyInfo_t
	cPolicy.version = C.dcgmPolicyInfo_version
	result := C.dcgmPolicyGetOne(handle.handle, C.ushort(fieldID), C.uint64_t(policyID), &cPolicy)
	if err := fieldPolicyResult(symbol, result); err != nil {
		return FieldPolicy{}, err
	}
	return fieldPolicyFromC(cPolicy), nil
}

// GetAllFieldPolicies returns the complete current field-policy list.
func GetAllFieldPolicies() ([]FieldPolicy, error) {
	const symbol = "dcgmPolicyGetAll"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return nil, err
	}

	var count C.uint
	result := C.dcgmPolicyGetAll(handle.handle, nil, &count)
	if result != C.DCGM_ST_OK && result != C.DCGM_ST_INSUFFICIENT_SIZE {
		return nil, fieldPolicyResult(symbol, result)
	}
	if count == 0 {
		return nil, nil
	}
	if err := validateFieldPolicyCount(int(count)); err != nil {
		return nil, err
	}

	cPolicies := make([]C.dcgmPolicyInfo_t, int(count))
	for i := range cPolicies {
		cPolicies[i].version = C.dcgmPolicyInfo_version
	}
	capacity := count
	result = C.dcgmPolicyGetAll(handle.handle, &cPolicies[0], &capacity)
	if err := fieldPolicyResult(symbol, result); err != nil {
		return nil, err
	}
	if int(capacity) > len(cPolicies) {
		return nil, fmt.Errorf("DCGM field-policy list grew from %d to %d entries; retry the request", len(cPolicies), capacity)
	}

	policies := make([]FieldPolicy, int(capacity))
	for i := range policies {
		policies[i] = fieldPolicyFromC(cPolicies[i])
	}
	return policies, nil
}

// ImportFieldPolicies imports field-policies from a DCGM YAML policy file.
func ImportFieldPolicies(path string) error {
	const symbol = "dcgmPolicyImport"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return err
	}
	if err := validateCString("field policy path", path, int(C.DCGM_POLICY_MAX_PATH)+1); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer freeCString(cPath)
	return fieldPolicyResult(symbol, C.dcgmPolicyImport(handle.handle, cPath))
}

// GetFieldPolicyViolations returns recent field-policy violations at or after since.
func GetFieldPolicyViolations(since time.Time) ([]FieldPolicyViolation, error) {
	const symbol = "dcgmPolicyGetViolations"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return nil, err
	}

	cViolations := make([]C.dcgm_policy_violation_t, MaxFieldPolicyViolations)
	count := C.uint(len(cViolations))
	result := C.dcgmPolicyGetViolations(handle.handle, C.int64_t(microsecondsOrZero(since)), &cViolations[0], &count)
	if result != C.DCGM_ST_OK && result != C.DCGM_ST_INSUFFICIENT_SIZE {
		return nil, fieldPolicyResult(symbol, result)
	}
	if int(count) > len(cViolations) {
		return nil, fmt.Errorf("DCGM reported %d policy violations; maximum supported is %d", count, len(cViolations))
	}

	violations := make([]FieldPolicyViolation, int(count))
	for i := range violations {
		violations[i] = fieldPolicyViolationFromC(cViolations[i])
	}
	return violations, nil
}

func boolToUint(value bool) uint {
	if value {
		return 1
	}
	return 0
}

func cFieldPolicyInfo(policy FieldPolicy) (C.dcgmPolicyInfo_t, error) {
	spec := FieldPolicySpec{
		FieldID:          policy.FieldID,
		Name:             policy.Name,
		Threshold:        policy.Threshold,
		Operator:         policy.Operator,
		Entities:         policy.Entities,
		Channels:         policy.Channels,
		Enabled:          policy.Enabled,
		RateLimitSeconds: policy.RateLimitSeconds,
	}
	if err := validateFieldPolicySpec(spec); err != nil {
		return C.dcgmPolicyInfo_t{}, err
	}

	var cPolicy C.dcgmPolicyInfo_t
	cPolicy.version = C.dcgmPolicyInfo_version
	cPolicy.policyId = C.uint64_t(policy.PolicyID)
	cPolicy.fieldId = C.ushort(policy.FieldID)
	copyGoStringToCChars(cPolicy.name[:], policy.Name)
	cPolicy.currentValue = C.double(policy.CurrentValue)
	cPolicy.threshold = C.double(policy.Threshold)
	cPolicy.policyOperator = C.dcgmPolicyOperator_t(policy.Operator)
	for i, entity := range policy.Entities {
		cPolicy.entities[i].entityGroupId = C.dcgm_field_entity_group_t(entity.EntityGroupId)
		cPolicy.entities[i].entityId = C.uint(entity.EntityId)
	}
	cPolicy.entityCount = C.uint(len(policy.Entities))
	cPolicy.channelMask = C.uint(policy.Channels)
	if policy.Enabled {
		cPolicy.state = C.DCGM_POLICY_STATE_ENABLED
	} else {
		cPolicy.state = C.DCGM_POLICY_STATE_DISABLED
	}
	cPolicy.rateLimitSec = C.uint(policy.RateLimitSeconds)
	cPolicy.createdAt = C.int64_t(microsecondsOrZero(policy.CreatedAt))
	cPolicy.violationCount = C.uint(policy.ViolationCount)
	cPolicy.lastTriggeredTimestamp = C.int64_t(microsecondsOrZero(policy.LastTriggeredTimestamp))
	return cPolicy, nil
}

func fieldPolicyFromC(cPolicy C.dcgmPolicyInfo_t) FieldPolicy {
	entityCount := min(int(cPolicy.entityCount), len(cPolicy.entities))
	entities := make([]GroupEntityPair, entityCount)
	for i := range entities {
		entities[i] = groupEntityPair(cPolicy.entities[i])
	}
	return FieldPolicy{
		PolicyID:               uint64(cPolicy.policyId),
		FieldID:                Short(cPolicy.fieldId),
		Name:                   C.GoString(&cPolicy.name[0]),
		CurrentValue:           float64(cPolicy.currentValue),
		Threshold:              float64(cPolicy.threshold),
		Operator:               FieldPolicyOperator(cPolicy.policyOperator),
		Entities:               entities,
		Channels:               FieldPolicyChannels(cPolicy.channelMask),
		Enabled:                cPolicy.state == C.DCGM_POLICY_STATE_ENABLED,
		RateLimitSeconds:       uint(cPolicy.rateLimitSec),
		CreatedAt:              timeFromMicroseconds(int64(cPolicy.createdAt)),
		ViolationCount:         uint(cPolicy.violationCount),
		LastTriggeredTimestamp: timeFromMicroseconds(int64(cPolicy.lastTriggeredTimestamp)),
	}
}

func fieldPolicyViolationFromC(cViolation C.dcgm_policy_violation_t) FieldPolicyViolation {
	return FieldPolicyViolation{
		PolicyID:     uint64(cViolation.policyId),
		PolicyName:   C.GoString(&cViolation.policyName[0]),
		FieldName:    C.GoString(&cViolation.fieldName[0]),
		FieldID:      Short(cViolation.fieldId),
		Entity:       groupEntityPair(cViolation.entity),
		CurrentValue: float64(cViolation.currentValue),
		Threshold:    float64(cViolation.thresholdValue),
		Operator:     FieldPolicyOperator(cViolation.policyOperator),
		Channels:     FieldPolicyChannels(cViolation.channelMask),
		Timestamp:    timeFromMicroseconds(int64(cViolation.timestamp)),
	}
}

func timeFromMicroseconds(value int64) time.Time {
	if value == 0 {
		return time.Time{}
	}
	return time.UnixMicro(value)
}

func microsecondsOrZero(value time.Time) int64 {
	if value.IsZero() {
		return 0
	}
	return value.UnixMicro()
}

func validateFieldPolicyCount(count int) error {
	if count > MaxFieldPolicies {
		return fmt.Errorf("DCGM reported %d field policies; maximum supported is %d", count, MaxFieldPolicies)
	}
	return nil
}

func copyGoStringToCChars(destination []C.char, value string) {
	for i := 0; i < len(value); i++ {
		destination[i] = C.char(value[i])
	}
}

type fieldPolicySubscriptionKey struct {
	fieldID  Short
	policyID uint64
}

type fieldPolicySubscriptionRegistry struct {
	mu            sync.Mutex
	nextID        uint64
	subscriptions map[uint64]*FieldPolicySubscription
	byKey         map[fieldPolicySubscriptionKey]uint64
}

func newFieldPolicySubscriptionRegistry() *fieldPolicySubscriptionRegistry {
	return &fieldPolicySubscriptionRegistry{
		subscriptions: make(map[uint64]*FieldPolicySubscription),
		byKey:         make(map[fieldPolicySubscriptionKey]uint64),
	}
}

// FieldPolicySubscription owns one DCGM v3 callback registration.
// Close unregisters the callback and closes Violations. Delivery is best-effort: a full channel drops a violation.
type FieldPolicySubscription struct {
	Violations <-chan FieldPolicyViolation

	registry *fieldPolicySubscriptionRegistry
	id       uint64
	key      fieldPolicySubscriptionKey
	ch       chan FieldPolicyViolation
	closeMu  sync.Mutex
	closed   bool
}

func (r *fieldPolicySubscriptionRegistry) add(key fieldPolicySubscriptionKey, buffer int) (*FieldPolicySubscription, error) {
	if buffer < 1 {
		buffer = 1
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byKey[key]; exists {
		return nil, fmt.Errorf("field-policy callback is already registered for field %d policy %d", key.fieldID, key.policyID)
	}
	r.nextID++
	if r.nextID == 0 {
		r.nextID++
	}
	channel := make(chan FieldPolicyViolation, buffer)
	subscription := &FieldPolicySubscription{
		Violations: channel,
		registry:   r,
		id:         r.nextID,
		key:        key,
		ch:         channel,
	}
	r.subscriptions[subscription.id] = subscription
	r.byKey[key] = subscription.id
	return subscription, nil
}

func (r *fieldPolicySubscriptionRegistry) remove(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	subscription, exists := r.subscriptions[id]
	if !exists {
		return
	}
	delete(r.subscriptions, id)
	delete(r.byKey, subscription.key)
	close(subscription.ch)
}

func (r *fieldPolicySubscriptionRegistry) deliver(id uint64, violation FieldPolicyViolation) {
	r.mu.Lock()
	defer r.mu.Unlock()
	subscription, exists := r.subscriptions[id]
	if !exists {
		return
	}
	select {
	case subscription.ch <- violation:
	default:
	}
}

var fieldPolicyCallbacks = newFieldPolicySubscriptionRegistry()

// SubscribeFieldPolicyViolations registers one DCGM v3 callback. A fieldID or policyID of zero is a DCGM wildcard.
func SubscribeFieldPolicyViolations(fieldID Short, policyID uint64, buffer int) (*FieldPolicySubscription, error) {
	const symbol = "dcgmPolicyRegister_v3"
	if err := requireFieldPolicySymbol(symbol); err != nil {
		return nil, err
	}

	subscription, err := fieldPolicyCallbacks.add(fieldPolicySubscriptionKey{fieldID: fieldID, policyID: policyID}, buffer)
	if err != nil {
		return nil, err
	}
	result := C.dcgmPolicyRegister_v3(
		handle.handle,
		C.ushort(fieldID),
		C.uint64_t(policyID),
		C.fpRecvPolicyViolation(C.fieldPolicyViolationNotify),
		C.uint64_t(subscription.id),
	)
	if err := fieldPolicyResult(symbol, result); err != nil {
		fieldPolicyCallbacks.remove(subscription.id)
		return nil, err
	}
	return subscription, nil
}

// Close unregisters this DCGM v3 callback. It is safe to call more than once.
func (s *FieldPolicySubscription) Close() error {
	return s.closeWith(func() error {
		const symbol = "dcgmPolicyUnregister_v3"
		if err := requireFieldPolicySymbol(symbol); err != nil {
			return err
		}
		result := C.dcgmPolicyUnregister_v3(handle.handle, C.ushort(s.key.fieldID), C.uint64_t(s.key.policyID))
		return fieldPolicyResult(symbol, result)
	})
}

func (s *FieldPolicySubscription) closeWith(unregister func() error) error {
	if s == nil {
		return nil
	}
	s.closeMu.Lock()
	defer s.closeMu.Unlock()
	if s.closed {
		return nil
	}

	if err := unregister(); err != nil && !unregisterErrorClearsLocalState(err) {
		return err
	}
	s.registry.remove(s.id)
	s.closed = true
	return nil
}

// FieldPolicyViolationRegistration copies a DCGM callback event into Go and delivers it without blocking the DCGM thread.
//
//export FieldPolicyViolationRegistration
func FieldPolicyViolationRegistration(data unsafe.Pointer, userData C.uint64_t) C.int {
	if data == nil {
		return 0
	}
	violation := *(*C.dcgm_policy_violation_t)(data)
	fieldPolicyCallbacks.deliver(uint64(userData), fieldPolicyViolationFromC(violation))
	return 0
}
