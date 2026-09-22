<div align="center">

# go-dcgm

**Go bindings for NVIDIA Data Center GPU Manager (DCGM).**

<a href="#development">Build and test</a>
·
<a href="#development-container">Development container</a>
·
<a href="samples/README.md">Samples</a>
·
<a href="CONTRIBUTING.md">Contribute</a>

<p>
  <a href="https://github.com/NVIDIA/go-dcgm/actions/workflows/go.yml"><img alt="GitHub Actions" src="https://github.com/NVIDIA/go-dcgm/actions/workflows/go.yml/badge.svg" /></a>
  <a href="https://pkg.go.dev/github.com/NVIDIA/go-dcgm"><img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/NVIDIA/go-dcgm.svg" /></a>
  <a href="LICENSE"><img alt="License: Apache-2.0" src="https://img.shields.io/badge/License-Apache--2.0-blue.svg" /></a>
  <a href="CONTRIBUTING.md"><img alt="DCO" src="https://img.shields.io/badge/DCO-1.1-blue.svg" /></a>
</p>

</div>

---

## Overview

This repository provides Go bindings for [NVIDIA Data Center GPU Manager
(DCGM)](https://developer.nvidia.com/dcgm). DCGM manages and monitors NVIDIA
GPUs in cluster environments.

The repository also includes sample programs that show how to use the bindings.

## Diagnostic API note

`RunDiag` returns hierarchical `DiagResults`. Existing `DiagResults.Software`
and `DiagResult` users remain supported through a deprecated compatibility
view. New code should use `results.Tests`. Each test includes its status,
per-entity results, errors, information, and extra data. See the [diagnostic
sample](samples/diag/main.go).

## Development

Use Task for common commands. Bazel builds the Go and CGO code.

These checks do not need a GPU:

```bash
task build
task test
task validate
```

Run GPU and DCGM tests only on a suitable system:

```bash
task test:integration
task test:race
```

### Development container

Open the repository with a [Dev Containers](https://containers.dev/) compatible
editor to use the checked-in development environment. It includes DCGM headers,
the pinned Go tools, and Bazel. CI builds and uses the same environment.

The normal build and validation commands do not need a GPU. GPU tests still
need a compatible NVIDIA host.

### Generate field constants

The DCGM field constants in `pkg/dcgm/const_fields.go` are automatically generated from `pkg/dcgm/dcgm_fields.h`. Curated lowercase compatibility names are tracked in `pkg/dcgm/legacy_fields.csv` and included during generation.

To regenerate these constants after updating the header file:

```bash
task generate
```

To verify that the generated code is up to date:

```bash
task generate:check
```

See [CONTRIBUTING.md](CONTRIBUTING.md#updating-dcgm-fields) for the full field
update process.

## Contribute and get help

Read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.

- Report a problem by [opening an issue](https://github.com/NVIDIA/go-dcgm/issues/new).
- Send a change by [opening a pull request](https://github.com/NVIDIA/go-dcgm).
