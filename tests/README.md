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

### Run All Tests

```bash
go test ./tests/... -v
```

### Run Specific Test Files

```bash
# Run device information tests
go test ./tests/deviceinfo_test.go -v

# Run monitoring tests
go test ./tests/dmon_test.go -v

# Run diagnostic tests
go test ./tests/diag_test.go -v
```

### Run Tests with Different Modes

```bash
# Run only quick tests (skip long-running tests)
go test ./tests/... -v -short

# Run tests with timeout
go test ./tests/... -v -timeout 5m
```

### Run Specific Test Functions

```bash
# Run specific test function
go test ./tests/deviceinfo_test.go -v -run TestDeviceInfo

# Run all tests matching a pattern
go test ./tests/... -v -run "TestDevice.*"
```

## Test Features

### Adaptive Testing

- Tests automatically skip when no GPUs are available
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
- Go 1.19+ for testing framework features

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
- Can run with or without actual GPU hardware (with appropriate skipping)

### Example GitHub Actions Integration

```yaml
- name: Run DCGM Tests
  run: |
    go test ./tests/... -v -timeout 10m
  continue-on-error: true  # Optional: allow failure if no GPU available
```

## Troubleshooting

### Common Issues

1. **No GPUs Found** - Tests will skip automatically
2. **Permission Denied** - Some tests require root privileges
3. **DCGM Not Available** - Ensure DCGM libraries are installed
4. **Timeout Issues** - Increase test timeout for slow systems

### Debug Information

All tests provide verbose logging when run with `-v` flag:

```bash
go test ./tests/deviceinfo_test.go -v
```

### Environment Variables

Tests respect standard Go testing environment variables:

- `GO_TEST_TIMEOUT_SCALE` - Scale test timeouts
- `DCGM_TESTING_MODE` - Custom testing configurations (if implemented)

## Contributing

When adding new tests:

1. Follow the existing naming pattern (`*_test.go`)
2. Include comprehensive documentation
3. Add appropriate test skipping for missing hardware
4. Include both positive and negative test cases
5. Update this README with new test descriptions
