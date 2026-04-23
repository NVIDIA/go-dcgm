/*
Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package dcgm provides Go bindings for the NVIDIA Data Center GPU Manager
// (DCGM) C library.
//
// # Overview
//
// DCGM is a suite of tools for managing and monitoring NVIDIA data-center GPUs
// in cluster environments. This package wraps the DCGM C API and exposes it as
// idiomatic Go, including health checks, field watches, diagnostics, GPU group
// management, policy violation monitoring, and topology queries.
//
// # Initialization
//
// Every program must call [Init] before using any other function, and should
// call the returned cleanup function (or [Shutdown]) before exiting.  Three
// operating modes are supported:
//
//   - [Embedded] – start the DCGM host engine inside the current process.
//     Suitable for standalone tools and tests.
//   - [Standalone] – connect to an already-running nv-hostengine daemon.
//     Pass the daemon's address as an additional argument.
//   - [StartHostengine] – spawn nv-hostengine as a child process, connect to
//     it, and automatically terminate it on shutdown.
//
// Example – embedded mode:
//
//	cleanup, err := dcgm.Init(dcgm.Embedded)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cleanup()
//
// Example – standalone mode (connect to a running nv-hostengine):
//
//	cleanup, err := dcgm.Init(dcgm.Standalone, "localhost:5555")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer cleanup()
//
// # Field Watches
//
// Monitoring GPU metrics is a three-step process:
//
//  1. Create a field group with the metric IDs you want to watch.
//  2. Create a GPU group (or use [GroupAllGPUs]) and start the watch.
//  3. Read values with [GetValuesSince] and clean up when done.
//
// Example:
//
//	cleanup, _ := dcgm.Init(dcgm.Embedded)
//	defer cleanup()
//
//	fields := []dcgm.Short{
//	    dcgm.DCGM_FI_DEV_GPU_TEMP,
//	    dcgm.DCGM_FI_DEV_POWER_USAGE,
//	}
//
//	fieldGroup, err := dcgm.FieldGroupCreate("myFields", fields)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dcgm.FieldGroupDestroy(fieldGroup)
//
//	gpuGroup := dcgm.GroupAllGPUs()
//	if err := dcgm.WatchFieldsWithGroup(fieldGroup, gpuGroup); err != nil {
//	    log.Fatal(err)
//	}
//	defer dcgm.UnwatchFields(fieldGroup, gpuGroup)
//
//	values, _, err := dcgm.GetValuesSince(gpuGroup, fieldGroup, time.Time{})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, v := range values {
//	    fmt.Printf("GPU %d field %d: %v\n", v.EntityID, v.FieldID, v.Int64())
//	}
//
// # GPU Groups
//
// GPU groups let you apply operations to a named set of GPUs.  Use
// [GroupAllGPUs] to target every GPU on the system, or [CreateGroup] to build
// a custom group.  Groups must be destroyed with [DestroyGroup] when no longer
// needed.
//
//	group, err := dcgm.CreateGroup("workers")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer dcgm.DestroyGroup(group)
//
//	// Add GPU 0 and GPU 1 to the group.
//	_ = dcgm.AddToGroup(group, 0)
//	_ = dcgm.AddToGroup(group, 1)
//
// # Health Checks
//
// Passive health monitoring tracks PCIe errors, NVLink faults, memory
// failures, and more.  Enable the watches for a group, then call [HealthCheck]
// to retrieve the current status.
//
//	dcgm.HealthSet(group, dcgm.DCGM_HEALTH_WATCH_ALL)
//	response, err := dcgm.HealthCheck(group)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, incident := range response.Incidents {
//	    fmt.Printf("GPU %d: %s – %s\n",
//	        incident.EntityInfo.EntityId,
//	        incident.System,
//	        incident.Health)
//	}
//
// # Policy Violation Monitoring
//
// Register callbacks for GPU policy violations (ECC errors, XID events, power
// limits, etc.) using [ListenForPolicyViolations].  The function returns a
// channel that receives [PolicyViolation] values.  Always cancel the context
// when monitoring is no longer required to avoid goroutine leaks.
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	violations, err := dcgm.ListenForPolicyViolations(ctx,
//	    dcgm.POLICY_COND_DBE,
//	    dcgm.POLICY_COND_NVLINK,
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for v := range violations {
//	    fmt.Printf("Policy violation on GPU %d: condition %v\n",
//	        v.Condition, v.Val)
//	}
//
// # Diagnostics
//
// Run the built-in GPU diagnostic suite with [RunDiag].  The diagnostics are
// graduated by level ([DiagQuick], [DiagShort], [DiagMedium], [DiagLong]).
//
//	results, err := dcgm.RunDiag(dcgm.DiagShort, dcgm.GroupAllGPUs())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, r := range results.Results {
//	    fmt.Printf("GPU %d test %s: %s\n", r.GPU, r.TestName, r.Result)
//	}
//
// # Thread Safety
//
// [Init] and [Shutdown] are protected by an internal mutex and are safe to
// call concurrently.  All other functions assume that [Init] has completed
// successfully.  The DCGM C library itself is thread-safe for read operations;
// consult the DCGM documentation for write-operation constraints.
//
// # Resource Management
//
// Many objects allocated by DCGM (groups, field groups, etc.) must be
// explicitly released.  The idiomatic pattern is to pair each Create/Watch
// call with a deferred Destroy/Unwatch call:
//
//	group, _ := dcgm.CreateGroup("g")
//	defer dcgm.DestroyGroup(group)
//
//	fg, _ := dcgm.FieldGroupCreate("fg", fields)
//	defer dcgm.FieldGroupDestroy(fg)
//
// Failing to release resources will cause memory leaks inside nv-hostengine
// for the duration of the process.
package dcgm
