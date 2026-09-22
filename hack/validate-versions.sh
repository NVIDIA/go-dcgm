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

expect_literal go.mod "go ${GO_VERSION}"
expect_literal .devcontainer/Dockerfile "COPY hack/versions.env"
expect_literal .github/workflows/go.yml "go-version: ${GO_VERSION}"
expect_literal .pre-commit-config.yaml "golang.org/x/tools/cmd/goimports@v${GOIMPORTS_VERSION}"
expect_literal Taskfile.yml "hack/versions.env"

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
