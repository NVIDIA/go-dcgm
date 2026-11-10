# hack - Agent Instructions

The root `AGENTS.md` still applies. This file holds the rules for shell
scripts anywhere in the repository.

## Shell scripts

- Start a new script from `script-template.sh`: `#!/usr/bin/env bash`,
  `set -euo pipefail`, `-h`/`--help`, usage and errors on standard error,
  flags read in one `case` loop, and the work in `main`. Indent with 2
  spaces, as the template does.
- Exceptions:
  - A script that takes only positional arguments needs strict mode and a
    usage line.
  - A script that a task runs with environment variables only, and no
    arguments, needs strict mode and must reject arguments with a usage
    line. It does not need the flag loop.
  - A script that must run under POSIX `sh` keeps `#!/bin/sh` and `set -eu`.
  - A sourced library has no `#!` line and starts with
    `# shellcheck shell=bash`. It does not set shell options for the script
    that sources it.
- Existing scripts adopt the pattern when you change them; do not rewrite
  scripts you do not otherwise touch.
- `task lint:shell` runs ShellCheck on the paths and globs listed in the
  task in `Taskfile.yml`. When you add a script outside those paths, add it
  to the list. Review strict-mode setup and any temporary option changes
  when editing a script.
