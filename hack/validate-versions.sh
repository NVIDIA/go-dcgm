#!/usr/bin/env bash

set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# shellcheck disable=SC1091
source "${root}/hack/versions.env"

status=0

expect_literal() {
  local file="$1"
  local literal="$2"

  if ! grep -Fq -- "${literal}" "${root}/${file}"; then
    printf '%s: expected %s\n' "${file}" "${literal}" >&2
    status=1
  fi
}

# A line that names a DCGM package version, gives DCGM_VERSION or
# DCGM_PACKAGE_REVISION a default in a Dockerfile ARG, or sets either one to
# a number (a shell, ENV, YAML or ${VAR:-default} assignment) copies the pin
# instead of reading it.
dcgm_copy_pattern='(datacenter-gpu-manager-4[a-z0-9-]*|Pin: version)[=[:space:]]+1:[0-9]|ARG[[:space:]]+DCGM_(VERSION|PACKAGE_REVISION)=|DCGM_(VERSION|PACKAGE_REVISION)[=:[:space:]-]+["'\'']?[0-9]'

# dcgm_copies prints each such line as file:line:text, from every tracked
# file but hack/versions.env in a git checkout, or from the development
# image's Dockerfile elsewhere.
dcgm_copies() {
  if git -C "${root}" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    git -C "${root}" grep -nE "${dcgm_copy_pattern}" -- ':!hack/versions.env' || true
  else
    grep -nE "${dcgm_copy_pattern}" "${root}/.devcontainer/Dockerfile" |
      sed 's|^|.devcontainer/Dockerfile:|' || true
  fi
}

main() {
  local number copy

  if (($# != 0)); then
    echo "Usage: ${0}" >&2
    exit 2
  fi

  expect_literal go.mod "go ${GO_VERSION}"
  expect_literal .devcontainer/Dockerfile "COPY hack/versions.env"
  expect_literal .devcontainer/Dockerfile "gocover-cobertura@v\${GOCOVER_COBERTURA_VERSION}"
  expect_literal .pre-commit-config.yaml "golang.org/x/tools/cmd/goimports@v${GOIMPORTS_VERSION}"
  expect_literal Taskfile.yml "hack/versions.env"

  # A number without leading zeros, so each pin has one release tag encoding.
  number='(0|[1-9][0-9]*)'
  if ! [[ "${DCGM_VERSION:-}" =~ ^${number}\.${number}\.${number}$ ]]; then
    printf 'hack/versions.env: DCGM_VERSION must be X.Y.Z without leading zeros, got %s\n' "${DCGM_VERSION:-}" >&2
    status=1
  fi
  if ! [[ "${DCGM_PACKAGE_REVISION:-}" =~ ^${number}$ ]]; then
    printf 'hack/versions.env: DCGM_PACKAGE_REVISION must be a number without leading zeros, got %s\n' "${DCGM_PACKAGE_REVISION:-}" >&2
    status=1
  fi
  expect_literal .devcontainer/Dockerfile "1:\${DCGM_VERSION}-\${DCGM_PACKAGE_REVISION}"

  while IFS= read -r copy; do
    printf '%s: names a DCGM version; use DCGM_VERSION and DCGM_PACKAGE_REVISION from hack/versions.env\n' \
      "${copy}" >&2
    status=1
  done < <(dcgm_copies)

  if [[ -f "${root}/.bazelversion" ]]; then
    expect_literal .bazelversion "${BAZEL_VERSION}"
  fi

  if [[ -f "${root}/MODULE.bazel" ]]; then
    expect_literal MODULE.bazel "bazel_dep(name = \"rules_go\", version = \"${RULES_GO_VERSION}\")"
    expect_literal MODULE.bazel "bazel_dep(name = \"gazelle\", version = \"${GAZELLE_VERSION}\")"
    expect_literal MODULE.bazel "bazel_dep(name = \"rules_cc\", version = \"${RULES_CC_VERSION}\")"
    expect_literal MODULE.bazel "go_sdk.download(version = \"${GO_VERSION}\")"
  fi

  exit "${status}"
}

main "$@"
