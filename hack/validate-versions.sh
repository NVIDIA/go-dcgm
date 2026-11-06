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
expect_literal .devcontainer/Dockerfile "gocover-cobertura@v\${GOCOVER_COBERTURA_VERSION}"
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

# active_value prints the value of the uncommented assignment "$2: value" in
# file $1, without surrounding quotes, or nothing when there is none.
active_value() {
    awk -v name="$2" '
        /^[[:space:]]*#/ { next }
        {
            line = $0
            sub(/^[[:space:]]+/, "", line)
            if (index(line, name ":") != 1) { next }
            value = substr(line, length(name) + 2)
            sub(/^[[:space:]]+/, "", value)
            sub(/[[:space:]]+(#.*)?$/, "", value)
            gsub(/^["\047]|["\047]$/, "", value)
            print value
            exit
        }' "$1"
}

# A registry variable reference, ${VAR} or $VAR, followed by a slash.
registry_prefix='^(\$\{[A-Za-z_][A-Za-z0-9_]*\}|\$[A-Za-z_][A-Za-z0-9_]*)/'

# check_image_variables checks the file $1, which sets CI image variables. A CI
# image variable is <NAME>_IMAGE: "${VAR}/<path>:<tag>@<digest>", where VAR is
# any variable that names a registry host. Each Docker Hub pin,
# <NAME>_IMAGE_REPOSITORY=docker.io/<path> with <NAME>_IMAGE_TAG and
# <NAME>_IMAGE_DIGEST, must be set there, uncommented, in that form, and each
# such variable in the file must have a Docker Hub pin.
check_image_variables() {
    local file="$1"
    local key name repository tag digest image actual variable

    while read -r key; do
        name="${key%_REPOSITORY}"
        repository="${!key}"
        tag="${name}_TAG"
        digest="${name}_DIGEST"

        if [[ "${repository}" != docker.io/* ]]; then
            continue
        fi

        if [[ -z "${!tag:-}" || -z "${!digest:-}" ]]; then
            printf 'hack/versions.env: %s needs %s and %s\n' "${key}" "${tag}" "${digest}" >&2
            status=1
            continue
        fi

        # Compare the part after the registry variable with the pin.
        image="${repository#docker.io/}:${!tag}@${!digest}"
        actual="$(active_value "${root}/${file}" "${name}")"
        if [[ -z "${actual}" ]]; then
            printf '%s: %s is not set; expected a registry variable followed by /%s\n' "${file}" "${name}" "${image}" >&2
            status=1
        elif ! [[ "${actual}" =~ ${registry_prefix} ]]; then
            printf '%s: %s is %s; expected a registry variable followed by /%s\n' "${file}" "${name}" "${actual}" "${image}" >&2
            status=1
        elif [[ "${actual#"${BASH_REMATCH[0]}"}" != "${image}" ]]; then
            printf '%s: %s is %s; hack/versions.env pins %s\n' "${file}" "${name}" "${actual}" "${BASH_REMATCH[0]}${image}" >&2
            status=1
        fi
    done < <(grep -oE '^[A-Z0-9_]+_IMAGE_REPOSITORY' "${root}/hack/versions.env" || true)

    # Every active CI image variable, one whose value is a registry variable,
    # a path and an @<digest>, needs a Docker Hub pin. The sed pattern matches
    # a literal $ before the registry variable.
    # shellcheck disable=SC2016
    while read -r variable; do
        key="${variable}_REPOSITORY"
        if [[ -z "${!key:-}" ]]; then
            printf '%s: %s needs %s=docker.io/<path> in hack/versions.env\n' \
                "${file}" "${variable}" "${key}" >&2
            status=1
        elif [[ "${!key}" != docker.io/* ]]; then
            printf '%s: %s uses a registry variable, but %s is %s\n' \
                "${file}" "${variable}" "${key}" "${!key}" >&2
            status=1
        fi
    done < <(grep -vE '^[[:space:]]*#' "${root}/${file}" |
        sed -nE 's/^[[:space:]]*([A-Z0-9_]+_IMAGE):[[:space:]]*["'\'']?\$(\{[A-Za-z_][A-Za-z0-9_]*\}|[A-Za-z_][A-Za-z0-9_]*)\/[^[:space:]]*@[^[:space:]]+.*/\1/p')
}

# An optional argument names the CI file that sets the image variables.
if (($# > 0)) && [[ -n "$1" ]]; then
    check_image_variables "$1"
fi

exit "${status}"
