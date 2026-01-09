# DCGM Fields Generator

This tool generates Go constants from the DCGM C header file `dcgm_fields.h`.

## Overview

The generator parses `dcgm_fields.h` and extracts all DCGM field definitions (`DCGM_FI_*` constants), then generates a Go file with:

- Typed constants for each DCGM field
- Field name mappings for lookup by string name
- Helper functions (`GetFieldID`, `GetFieldIDOrPanic`, etc.)
- Legacy field mappings for backward compatibility

## Usage

The generator is typically invoked via `go generate` or `make generate`:

```bash
# Via Make
make generate

# Via go generate
go generate ./...
```

### Direct Usage

You can also run the generator directly:

```bash
go run cmd/gen-fields/main.go cmd/gen-fields/template.go \
    pkg/dcgm/dcgm_fields.h \
    pkg/dcgm/const_fields.go
```

Arguments:
1. Path to `dcgm_fields.h` (input)
2. Path to `const_fields.go` (output)

## How It Works

1. **Parse Header File**: Reads `dcgm_fields.h` and extracts all `#define DCGM_FI_*` definitions
2. **Extract Field Information**:
   - Field name (e.g., `DCGM_FI_DEV_GPU_TEMP`)
   - Field ID (numeric value)
   - Field comment/description
3. **Generate Go Code**: Uses Go templates to create:
   - Constant definitions: `DCGM_FI_DEV_GPU_TEMP Short = 150`
   - Field name maps for string-based lookup
   - Helper functions for field ID resolution

## Output

The generated `const_fields.go` file contains:

```go
const (
    DCGM_FI_DEV_GPU_TEMP Short = 150
    DCGM_FI_DEV_POWER_USAGE Short = 155
    // ... etc
)

var dcgmFields = map[string]Short{
    "dcgm_gpu_temp": 150,
    "dcgm_power_usage": 155,
    // ... etc
}

func GetFieldID(fieldName string) (Short, bool) { ... }
func GetFieldIDOrPanic(fieldName string) Short { ... }
```

## Template

The code generation template is defined in `template.go` and includes the full structure of the output Go file.

## Updating Fields

When DCGM adds new fields:

1. Update `pkg/dcgm/dcgm_fields.h` with the latest version from DCGM
2. Run `make generate`
3. Review the diff in `pkg/dcgm/const_fields.go`
4. Commit both the header and generated file

See [CONTRIBUTING.md](../../CONTRIBUTING.md#updating-dcgm-fields) for detailed instructions.

