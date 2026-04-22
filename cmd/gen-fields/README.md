# DCGM Fields Generator

This tool generates Go constants from the DCGM C header file `dcgm_fields.h`.

## Overview

The generator parses `dcgm_fields.h` and generates a Go file with:

- Typed constants for each `#define DCGM_FI_X <int>` in the header.
- `dcgmFields`: maps canonical name to field ID.
- `legacyDCGMFields`: maps backward-compatible names to the same IDs.
  Populated from two sources:
    - Hand-curated DCGM 1.x era lowercase names (e.g. `dcgm_gpu_temp`),
      preserved across regenerations.
    - Deprecated-alias `#define OLD NEW` lines in the header, either
      inside an `#ifdef DCGM_DEPRECATED` block or preceded by a
      `Deprecated:` comment.
- Helper functions: `GetFieldID`, `GetFieldIDOrPanic`, `IsLegacyField`,
  `IsCurrentField`.

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

1. **Parse header**: reads `dcgm_fields.h`. Two shapes of `#define` are extracted:
   - `#define DCGM_FI_X <int>` -> canonical field, with preceding comment
     captured as its description.
   - `#define DCGM_FI_OLD DCGM_FI_NEW` -> deprecated alias. Recorded only
     when the alias is inside an `#ifdef DCGM_DEPRECATED` block OR its
     preceding comment contains `Deprecated:`. Other alias-style
     `#define`s (e.g. range sentinels) are silently skipped.
2. **Resolve aliases**: each recorded alias is mapped to its target
   field's canonical ID. If a target isn't a known field, generation
   fails so header churn can't silently drop previously-exposed names.
3. **Preserve curated legacy names**: lowercase DCGM 1.x names in the
   previously-generated `const_fields.go` are round-tripped. `DCGM_FI_*`
   entries are not preserved here; they re-derive from step 2 every run.
4. **Emit Go code** via `template.go`, then run `gofmt -w` on the output
   so `make check-generate` stays stable.

## Output

The generated `const_fields.go` file contains:

```go
const (
    DCGM_FI_DEV_GPU_TEMP    Short = 150
    DCGM_FI_DEV_POWER_USAGE Short = 155
    // ... etc
)

var dcgmFields = map[string]Short{
    "DCGM_FI_DEV_GPU_TEMP":    150,
    "DCGM_FI_DEV_POWER_USAGE": 155,
    // ... etc
}

var legacyDCGMFields = map[string]Short{
    // DCGM 1.x lowercase names, preserved from the prior generation:
    "dcgm_gpu_temp":   150,
    "dcgm_power_usage": 155,
    // Deprecated aliases resolved from dcgm_fields.h:
    "DCGM_FI_DEV_CLOCK_THROTTLE_REASONS":  112,
    "DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL": 449,
    // ... etc
}

func GetFieldID(fieldName string) (Short, bool) { ... }
func GetFieldIDOrPanic(fieldName string) Short { ... }
func IsLegacyField(fieldName string) bool     { ... }
func IsCurrentField(fieldName string) bool    { ... }
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

