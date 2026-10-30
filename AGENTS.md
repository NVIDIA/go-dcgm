# go-dcgm - Agent Instructions

go-dcgm is the Go module `github.com/NVIDIA/go-dcgm`: CGO bindings for NVIDIA
Data Center GPU Manager (DCGM), plus sample programs that show how to use them.
Other Go projects import `pkg/dcgm`; this repository ships no binary.

## Source of truth

When sources disagree, trust them in this order:

1. Current Go code, tests, Bazel `BUILD.bazel` files, and checked-in headers.
2. `Taskfile.yml` and the CI workflow files.
3. `go.mod`, `MODULE.bazel`, and `hack/versions.env` for versions and pins.
4. The official DCGM documentation for API and field meaning.
5. `README.md`, `CONTRIBUTING.md`, folder READMEs, and `llms.txt`.

## Working rules

- Do not invent DCGM functions, field IDs, flags, task names, or Bazel
  targets. Find them in the code, headers, or `Taskfile.yml` first.
- Never report a check as passing when it did not run. Name skipped checks and
  why, such as no GPU, no DCGM, no Docker, or no network.
- If the same approach fails three times, stop and explain what you tried.
- Ask before changing an exported API. Keep old names working with a
  `// Deprecated:` comment, as `const.go` and `diag.go` do.
- Keep changes scoped. Do not commit local settings, `bazel-*` links, logs,
  coverage output, tool caches, or plans and their evidence.

## Repository map

| Path | What it owns |
| --- | --- |
| `pkg/dcgm/` | The bindings, copied DCGM C headers, and package tests |
| `cmd/gen-fields/` | Generates `pkg/dcgm/const_fields.go` from `dcgm_fields.h` |
| `samples/` | Example programs, including the `restApi` HTTP sample |
| `tests/` | GPU and DCGM integration tests that mirror the samples |
| `hack/` | Version pins (`versions.env`), pin check, CI coverage script |
| `.github/workflows/` | GitHub build, test, and lint workflow |
| `.devcontainer/` | Development image used locally and by CI |

`pkg/dcgm/dcgm_*.h` and `nvml.h` are copied from DCGM. Do not edit them by
hand; follow [Updating DCGM Fields](CONTRIBUTING.md#updating-dcgm-fields).

## Commands

Use the tasks in `Taskfile.yml`; run `task --list` to see them. Bazel builds
the Go and CGO code.

These need no GPU:

```sh
task build              # library and samples; checks generated code first
task test               # same as task test:unit
task test:unit          # prerequisite-free Bazel tests
task fmt:check          # gofumpt
task lint               # golangci-lint; lint:incremental checks since HEAD~1
task generate:check     # generated constants and mocks are current
task versions:validate  # version pins agree
task bazel:validate     # Gazelle, Bzlmod lockfile, and CI-safe targets
task validate           # pins, Bazel, generated code, format, and lint
task test:error-metadata # needs DCGM libraries, as in the dev container
```

These need a host with a GPU and the DCGM runtime:

```sh
task test:gpu:prereqs      # confirms DCGM finds a physical GPU
task test:integration      # integration suite; runs the prereq check first
task test:race             # integration suite with the race detector
task test:diagnostics      # hardware diagnostic tests
```

The test tasks (`test:unit`, `test:error-metadata`, `test:integration`,
`test:race`, `test:diagnostics`) pass Bazel options given after `--`, such as
`-- --test_filter=TestDeviceInfo`. `task build` ignores them. In
`test:integration` and `test:race`, the options reach only the main suite run,
not the first policy test run.
`task fmt` and `task generate` change files. `task bazel:gazelle` and
`task bazel:tidy` update Bazel files. Docs-only changes need
`git diff --check` and `pre-commit run --all-files` (markdownlint, yamllint).
Run containers that write into the checkout with
`--user "$(id -u):$(id -g)"`.

## Code and comments

- Reuse existing code, then the standard library, then existing
  dependencies. Ask before adding a dependency.
- Comments explain the code. Never mention plans, prompts, tickets, or why an
  agent made a change. No emojis.
- Do not edit `pkg/dcgm/const_fields.go` or `native_calls_mock_test.go` by
  hand. Run `task generate` after changing `dcgm_fields.h`,
  `legacy_fields.csv`, `cmd/gen-fields`, or `native_calls.go`.

### Go

- Use the Go version in `go.mod` and the tool versions in `hack/versions.env`.
- Return errors with context (`fmt.Errorf("...: %w", err)`). Do not drop
  errors from DCGM calls.
- Free C memory and DCGM handles you create on every path, including errors.
- Format with `task fmt`. Use table tests for contract cases.
- Every `//nolint` names the linter and gives a reason.

## Tests

- `pkg/dcgm` has two Bazel test targets. `dcgm_unit_test` needs no GPU and
  uses gomock fakes of the native calls. `dcgm_test` links the DCGM runtime.
  Add each new test file to the right target in `pkg/dcgm/BUILD.bazel`. Files
  for `dcgm_unit_test` also go in the `gazelle:exclude` list at its top.
- Everything under `tests/` needs a GPU and DCGM. Many tests call `t.Skip`
  when a GPU feature or root access is missing; report those as skipped.
- Test at the lowest layer that proves the behavior. Do not write unit or
  mock tests for CI files or scripts; run the real command or CI job.
- Changes to runtime behavior also need `task test:integration` and
  `task test:race` on a GPU host. If you could not run them, say so.

## Commits and reviews

- Use Conventional Commits with a scope, such as `fix(health): ...`,
  `test(dcgm): ...`, or `chore(deps): ...`. The subject says what changed.
- Sign off every commit with `git commit -s` (DCO; see `CONTRIBUTING.md`).
- Do not push to `main`. Open a pull request or merge request. Fix every
  supported review finding in scope, or reply with evidence; never resolve
  silently.

## Keeping this file current

- Add a rule only for a repeated mistake or a fact an agent cannot easily find.
  Prefer a check, linter rule, test, or task over a written rule.
- No status, dates, version numbers, or plan IDs. Keep this file under 200
  lines; put folder-only rules in that folder's `AGENTS.md`.
- When you change this file, remove rules that no longer apply.
- `llms.txt` is the usage contract for programs that call the bindings.
  Update it when an exported API it names changes.
