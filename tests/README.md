# DCGM Go Testing Samples

This directory contains test versions of all the DCGM samples, reimplemented using the Go testing framework. These tests demonstrate the functionality of the NVIDIA Data Center GPU Manager (DCGM) Go bindings while being suitable for automated testing and CI/CD pipelines.

## Test Files Overview

### Core Device Management

- **`deviceinfo_test.go`** - Tests device information retrieval functionality
  - Equivalent to `samples/deviceInfo/main.go`
  - Tests GPU device properties, identification, and topology information
  - Includes tests for both embedded and standalone hostengine connections

- **`dmon_test.go`** - Tests device monitoring capabilities
  - Equivalent to `samples/dmon/main.go`
  - Monitors GPU utilization, temperature, power, and clock speeds
  - Includes time-limited monitoring tests and sample consistency checks

- **`device_status_test.go`** - Tests device status querying (part of dmon functionality)
  - Tests single and multiple GPU status queries
  - Validates utilization metrics and system health indicators

### Diagnostics and Health

- **`diag_test.go`** - Tests DCGM diagnostic functionality
  - Equivalent to `samples/diag/main.go`
  - Runs quick and medium-level diagnostic tests
  - Validates software and hardware diagnostic results

- **`health_test.go`** - Tests GPU health monitoring
  - Equivalent to `samples/health/main.go`
  - Performs single and continuous health checks
  - Tests health watch configuration and error reporting

### System Management

- **`hostengine_test.go`** - Tests DCGM hostengine introspection
  - Equivalent to `samples/hostengineStatus/main.go`
  - Monitors hostengine memory and CPU usage
  - Tests introspection under different load conditions

- **`policy_test.go`** - Tests policy violation monitoring
  - Equivalent to `samples/policy/main.go`
  - Tests various policy condition types (DBE, XID, thermal, power)
  - Includes context cancellation and timeout handling

### Process and Topology

- **`processinfo_test.go`** - Tests GPU process monitoring
  - Equivalent to `samples/processInfo/main.go`
  - Tests process field watching and information retrieval
  - Includes PID-specific testing capabilities

- **`topology_test.go`** - Tests GPU topology analysis
  - Equivalent to `samples/topology/main.go`
  - Tests inter-GPU connection discovery and analysis
  - Includes topology consistency validation

### REST API

- **`restapi_test.go`** - Tests REST API endpoint functionality
  - Equivalent to `samples/restApi/` (complete implementation)
  - Uses `httptest` for testing HTTP endpoints without starting a real server
  - Tests JSON response formats and error handling

## Running the Tests

### Run Integration Tests

```bash
task test:integration
```

To run the full GPU checks on a qualified host:

```bash
task test:integration REPORTS=true
task test:race REPORTS=true
task test:diagnostics REPORTS=true
```

### Run Focused Tests

```bash
task test:integration -- --test_filter=TestDeviceInfo
task test:integration -- --test_filter=TestDeviceMonitoring
```

Integration and race run `//pkg/dcgm:dcgm_integration_test` before the
sample tests. Options after `--` apply to the sample tests. To focus on
a Go binding test, run its target directly.

```bash
bazel test --config=ci --config=dcgm-runtime --cache_test_results=no \
  //pkg/dcgm:dcgm_integration_test --test_filter='^TestPolicyErrors$'
```

`task test:diagnostics` runs `//pkg/dcgm:dcgm_diagnostics_test` and then the
sample diagnostic target (`TestDiagnostics`, `TestDiagnosticsLong`, and
`TestDiagTestNameFormat`). Options after `--` apply to both diagnostic runs.
To run one Go binding diagnostic test on a GPU host, use its Bazel target.

```bash
bazel test --config=ci --config=dcgm-runtime --cache_test_results=no \
  --test_env=DCGM_TEST_DIAGNOSTICS=1 //pkg/dcgm:dcgm_diagnostics_test \
  --test_filter='^TestDiagnosticNamedParameters$'
```

The broader `//pkg/dcgm:dcgm_test` target also includes tests that need
the DCGM library and run outside CI. Use `task test:policy:rediscovery` for the opt-in
policy driver detach and reattach check on a dedicated GPU host.

### Run Tests with Different Modes

```bash
# Run prerequisite-free unit tests
task test

# Run integration tests with a five-minute timeout
task test:integration -- --test_timeout=300
```

### Run Specific Test Functions

```bash
# Run specific test function
task test:integration -- --test_filter=TestDeviceInfo

# Run all tests matching a pattern
task test:integration -- --test_filter='TestDevice.*'
```

## Test Features

### Adaptive Testing

- The GPU test tasks fail when a GPU is unavailable. `task test:gpu:prereqs`
  checks that embedded DCGM starts and discovers a physical GPU before the
  integration, race, and diagnostic suites run.
- Different behavior for single vs. multi-GPU systems
- Graceful handling of permission-restricted operations

### Time-Limited Execution

- Long-running samples (like monitoring) are time-limited in tests
- Configurable test durations for CI/CD environments
- Background operations are properly cancelled

### Comprehensive Coverage

- Each test covers the core functionality of its corresponding sample
- Additional test scenarios for error conditions and edge cases
- Validation of return values and data consistency

### CI/CD Friendly

- Tests use the Go testing framework's standard patterns
- Proper test isolation and cleanup
- Structured logging for debugging

## Prerequisites

### System Requirements

- NVIDIA GPU(s) with DCGM support
- NVIDIA drivers installed
- DCGM libraries available
- The Go toolchain version pinned in [`hack/versions.env`](../hack/versions.env)

### Dependencies

The tests require the same dependencies as the original samples:

- `github.com/NVIDIA/go-dcgm/pkg/dcgm`
- `github.com/gorilla/mux` (for REST API tests only)

### Permissions

Some tests may require elevated privileges:

- Process monitoring tests work best when run as root
- Certain policy violation tests require administrative access
- Diagnostic tests may need elevated permissions for hardware access

## Test Structure

Each test file follows a consistent pattern:

1. **Basic Functionality Test** - Core sample functionality
2. **Extended Tests** - Additional scenarios and edge cases
3. **Error Handling Tests** - Validation of error conditions
4. **Performance/Consistency Tests** - Multi-sample validation

### Example Test Pattern

```go
func TestSampleFunctionality(t *testing.T) {
    // Initialize DCGM
    cleanup, err := dcgm.Init(dcgm.Embedded)
    if err != nil {
        t.Fatalf("Failed to initialize DCGM: %v", err)
    }
    defer cleanup()

    // Test core functionality
    // ... test implementation

    // Validate results
    // ... assertions and checks
}
```

## Integration with CI/CD

These tests are designed to integrate well with continuous integration systems:

- Use standard Go testing patterns
- Provide detailed logging for troubleshooting
- Support timeout and cancellation
- Run GPU test tasks on a qualified GPU host with the DCGM runtime. The
  tasks fail if they find no GPU, cannot load DCGM, or a test fails. The
  tasks run without cached results.
  Tests that need optional hardware or access, such as multiple GPUs, a CPU
  hierarchy, root access, a standalone hostengine, or the ConnectX diagnostic
  plugin, may skip and record why. The diagnostic sample tests require at
  least one passing result from DCGM and fail if DCGM reports any failed
  check. They log the details and errors that DCGM returns.

Test tasks share event and report paths within a checkout. Run them serially;
use separate worktrees for concurrent runs.

### Example GitHub Actions Integration

```yaml
- name: Run DCGM Tests
  run: task test:integration -- --test_timeout=600
```

`task test:gpu:prereqs`, which the GPU test tasks run first, runs `nvidia-smi`,
checks that `libdcgm.so.4` is installed, and checks that embedded DCGM starts
and discovers a physical GPU. None of these treats a missing GPU as a pass.

## Troubleshooting

### Common Issues

1. **No GPUs Found** - A GPU test task fails;
   run it on a qualified GPU/DCGM system
2. **Permission Denied** - Some tests require root privileges
3. **DCGM Not Available** - Ensure DCGM libraries are installed
4. **Timeout Issues** - Increase test timeout for slow systems

### Debug Information

Request complete Bazel test output when debugging:

```bash
task test:integration -- --test_output=all --test_filter=TestDeviceInfo
```

### Bazel Options

Focused Bazel options can be passed after `--`, including `--test_filter`,
`--test_timeout`, and `--test_output`. In `task test:diagnostics`, those
options apply to both the Go binding and sample tests. Omit them to run
the default tests.

## Contributing

When adding new tests:

1. Follow the existing naming pattern (`*_test.go`)
2. Include comprehensive documentation
3. Keep the GPU preflight failure for missing hardware; skip only tests for
   capabilities that are optional on an otherwise qualified system
4. Include both positive and negative test cases
5. Update this README with new test descriptions

### Job statistics with GPU load

`TestJobStats` requires root and a GPU-load executable. It uses a private
embedded engine, creates its own group and job, and removes both on every path.
It snapshots and restores GPU accounting mode using `nvidia-smi`.
It also checks `JobRemoveAll` within that engine. This case now runs through
`test:integration`.

```sh
task test:integration -- --test_filter='^TestJobStats$' --test_env=DCGM_JOB_LOAD_COMMAND=/absolute/path/to/gpu-load
```

The executable is invoked directly without arguments or a shell. It must read
`DCGM_JOB_GPU_ID` (the DCGM physical GPU ID), `DCGM_JOB_GPU_UUID` and
`DCGM_JOB_LOAD_SECONDS=3`, run
real compute work on that GPU for about three seconds, synchronize, and exit
successfully. The test kills it after 15 seconds and allows one second for
output pipes to close. The executable must not fork child processes. A small
CUDA program matching `DCGM_JOB_GPU_UUID` to a CUDA device and running a
timed kernel is suitable; no load generator is bundled or downloaded.
Without `DCGM_JOB_LOAD_COMMAND` the test reports a skip, not a lifecycle pass.

### Live PID query with GPU load

`TestProcessInfoWithGPUFixture` uses the same executable contract as `TestJobStats`,
but requests `DCGM_JOB_LOAD_SECONDS=10`. Run as root on a GPU with process
accounting support. It starts PID watches before launching the executable,
polls its live subprocess PID until DCGM returns a nonempty record for the
selected GPU, and checks the PID and legacy timestamps. It requires the
executable to exit successfully within 20 seconds, with one second for output
pipes to close. It restores accounting mode and destroys its owned group.

```sh
task test:integration -- --test_filter='^TestProcessInfoWithGPUFixture$' --test_env=DCGM_JOB_LOAD_COMMAND=/absolute/path/to/gpu-load
task test:race -- --test_filter='^TestProcessInfoWithGPUFixture$' --test_env=DCGM_JOB_LOAD_COMMAND=/absolute/path/to/gpu-load
```

Absent opt-in, root access, or PID/accounting capability is reported as a skip;
it does not qualify the PID runtime path. Native query errors or missing data
cannot pass this case. `TestJobStats` alone does not prove live PID queries.

### GPU configuration

`task test:configuration` opts in to the root-only configuration suite. Normal
integration and race tasks skip these checks. The suite uses a private embedded
engine and an owned single-GPU group. It snapshots both live and target state,
records and enforces the current power limit without changing its watts, clears
the target with `ConfigReset`, restores the original live and target settings,
and compares both snapshots. It skips before mutation if a preexisting target
would make enforcing or restoring it change unrelated settings.

The workload profile case reads metadata and profile names, then reapplies the
existing requested mask without selecting a different policy. It restores the
mask and target even on failure, and checks both live and target configuration.
A GPU or driver without workload profiles reports a separate capability skip;
that skip is not proof that profile configuration works.

```sh
task test:configuration
task test:configuration -- --config=race
```
