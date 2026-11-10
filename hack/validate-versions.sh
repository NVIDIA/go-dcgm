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

# check_image_variables checks the file $1, which sets CI image variables. A CI
# image variable is <NAME>_IMAGE: "${VAR}/<path>:<tag>@<digest>", where VAR is
# $2, the variable that names the registry the images are pulled through. Each
# Docker Hub pin, <NAME>_IMAGE_REPOSITORY=docker.io/<path> with
# <NAME>_IMAGE_TAG and <NAME>_IMAGE_DIGEST, must be set there, uncommented, in
# that form. check_registry_references checks the other direction.
check_image_variables() {
    local file="$1"
    local registry="$2"
    local key name repository tag digest image actual

    if ! [[ "${registry}" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
        printf '%s: IMAGE_REGISTRY_VARIABLE_NAME is "%s"; set it to the name of the variable the image values start with\n' \
            "${file}" "${registry}" >&2
        status=1
        return
    fi

    # ${VAR}/ or $VAR/ at the start of a value. \$ is a literal $, so nothing
    # here expands.
    local expected_prefix="\${${registry}}/"
    local registry_prefix='^(\$\{'"${registry}"'\}|\$'"${registry}"')/'

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
            printf '%s: %s is not set; expected %s%s\n' "${file}" "${name}" "${expected_prefix}" "${image}" >&2
            status=1
        elif ! [[ "${actual}" =~ ${registry_prefix} ]]; then
            printf '%s: %s is %s; expected %s%s\n' "${file}" "${name}" "${actual}" "${expected_prefix}" "${image}" >&2
            status=1
        elif [[ "${actual#"${BASH_REMATCH[0]}"}" != "${image}" ]]; then
            printf '%s: %s is %s; hack/versions.env pins %s\n' "${file}" "${name}" "${actual}" "${BASH_REMATCH[0]}${image}" >&2
            status=1
        fi
    done < <(grep -oE '^[A-Z0-9_]+_IMAGE_REPOSITORY' "${root}/hack/versions.env" || true)
}

# registry_values prints "line key image" for every reference to ${VAR}/ or
# $VAR/ in file $1, where VAR is $2. It reads each line in four steps:
#
#   1. Drop the comment, if any.
#   2. Find the key that starts the line, such as FOO_IMAGE in
#      "  - FOO_IMAGE: ...".
#   3. For each reference, left to right, take the key right before it in a
#      flow mapping, such as NAME in "{ NAME: ${VAR}/... }", or else the line
#      key. Text such as echo "NAME: ${VAR}/..." in a script line is not a
#      flow mapping, so it names no key.
#   4. Read the image after VAR/ up to the character that ends it.
#
# The fields are separated by the unit separator (octal 037), which read
# keeps apart even when the key is empty.
registry_values() {
    REGISTRY="$2" awk '
        BEGIN {
            braced = "${" ENVIRON["REGISTRY"] "}/"
            plain = "$" ENVIRON["REGISTRY"] "/"

            # Patterns for keys, built from smaller named parts. \047 is a
            # single quote, which cannot appear inside this shell-quoted
            # program.
            name = "[A-Za-z_][A-Za-z0-9_-]*"
            quote = "[\"\047]?"
            blanks = "[[:space:]]*"
            # "  - NAME:" or "  \"NAME\":" at the start of a line.
            line_key_pattern = "^" blanks "(-[[:space:]]+)?" quote name quote blanks ":"
            # "{ NAME: " or ", NAME: \"" at the end of the text before a reference.
            flow_key_pattern = "[{,]" blanks quote name quote blanks ":" blanks quote "$"

            # An image is literal characters and ${NAME} or $NAME parts. It ends
            # at a blank, a quote, a comma, a semicolon, ), ], a } that closes
            # no ${, or a $ that starts no variable.
            variable = "[A-Za-z_][A-Za-z0-9_]*"
            braced_variable = "^[$][{]" variable "[}]"
            plain_variable = "^[$]" variable
            image_ends = " \t\r\f\v\"\047,;)]}$"
        }

        # strip_comment returns line s without its comment. A # starts a
        # comment at the start of the line or after a space or tab, unless it
        # is inside quotes. A quote opens quoted text only where a quoted
        # value or word can start: at the start of the line or after a blank,
        # :, [, {, = or a comma, so an apostrophe inside a word opens nothing. As in
        # YAML, a backslash escapes the next character only inside double
        # quotes.
        function strip_comment(s,    i, c, in_quote) {
            in_quote = ""
            for (i = 1; i <= length(s); i++) {
                c = substr(s, i, 1)

                if (in_quote != "") {
                    if (in_quote == "\"" && c == "\\") { i++ }
                    else if (c == in_quote) { in_quote = "" }
                    continue
                }

                if ((c == "\"" || c == "\047") && (i == 1 || substr(s, i - 1, 1) ~ /[ \t:[{,=]/)) {
                    in_quote = c
                } else if (c == "#" && (i == 1 || substr(s, i - 1, 1) ~ /[ \t]/)) {
                    return substr(s, 1, i - 1)
                }
            }
            return s
        }

        # key_in returns the key name in the part of text that matches
        # pattern, or "" when nothing matches. Only the key itself in that
        # part looks like a name; the rest is blanks, quotes and punctuation.
        function key_in(text, pattern) {
            if (!match(text, pattern)) { return "" }
            text = substr(text, RSTART, RLENGTH)
            match(text, name)
            return substr(text, RSTART, RLENGTH)
        }

        # read_image returns the image at the start of s.
        function read_image(s,    image, c) {
            image = ""
            while (s != "") {
                if (match(s, braced_variable) || match(s, plain_variable)) {
                    image = image substr(s, 1, RLENGTH)
                    s = substr(s, RLENGTH + 1)
                    continue
                }

                c = substr(s, 1, 1)
                if (index(image_ends, c)) { break }
                image = image c
                s = substr(s, 2)
            }
            return image
        }

        {
            line = strip_comment($0)
            line_key = key_in(line, line_key_pattern)

            # rest is the part of the line after the last reference read.
            rest = line
            while (1) {
                b = index(rest, braced)
                p = index(rest, plain)
                if (b == 0 && p == 0) { break }

                # Take whichever form of the reference comes first.
                if (b != 0 && (p == 0 || b < p)) { start = b; prefix = braced }
                else { start = p; prefix = plain }

                key = key_in(substr(rest, 1, start - 1), flow_key_pattern)
                if (key == "") { key = line_key }

                rest = substr(rest, start + length(prefix))
                image = read_image(rest)
                rest = substr(rest, length(image) + 1)

                printf "%d\037%s\037%s\n", NR, key, image
            }
        }' "$1"
}

# digest_problem prints why image $1 is not pinned, or nothing when it ends
# in @sha256: and 64 hex digits, or in @${NAME} or @$NAME where NAME ends in
# _DIGEST and hack/versions.env sets it to such a digest. It reads that file
# under ${root}, and NAME's value from the environment the script sourced
# from it.
digest_problem() {
    local image="$1"
    local pinned='@(sha256:[0-9a-f]{64}|[$][{]([A-Za-z_][A-Za-z0-9_]*_DIGEST)[}]|[$]([A-Za-z_][A-Za-z0-9_]*_DIGEST))$'
    local name

    if ! [[ "${image}" =~ ${pinned} ]]; then
        echo "without an @sha256: digest or a <NAME>_DIGEST variable"
        return
    fi

    # A literal digest leaves no variable name.
    name="${BASH_REMATCH[2]}${BASH_REMATCH[3]}"
    [[ -n "${name}" ]] || return 0

    if ! grep -q "^${name}=" "${root}/hack/versions.env"; then
        echo "with ${name}, which hack/versions.env does not set"
        return
    fi

    if ! [[ "${!name}" =~ ^sha256:[0-9a-f]{64}$ ]]; then
        echo "with ${name}, which hack/versions.env sets to ${!name}, not sha256: and 64 hex digits"
    fi
}

# check_registry_references fails, for the files after $1, for each reference
# to the registry variable $1 that digest_problem rejects, and for each
# <NAME>_IMAGE that uses it without a Docker Hub pin in hack/versions.env.
check_registry_references() {
    local registry="$1"
    shift
    local file line key image pin problem
    local -A seen=()

    for file in "$@"; do
        # A glob can name the image variables file again.
        [[ -z "${seen[${file}]:-}" ]] || continue
        seen[${file}]=1

        if [[ ! -f "${root}/${file}" ]]; then
            printf '%s: no such file\n' "${file}" >&2
            status=1
            continue
        fi

        while IFS=$'\037' read -r line key image; do
            problem="$(digest_problem "${image}")"
            if [[ -n "${problem}" ]]; then
                printf '%s:%s: %s pulls through %s %s\n' \
                    "${file}" "${line}" "${image}" "${registry}" "${problem}" >&2
                status=1
            fi

            [[ "${key}" =~ ^[A-Z0-9_]+_IMAGE$ ]] || continue
            pin="${key}_REPOSITORY"
            if [[ -z "${!pin:-}" ]]; then
                printf '%s:%s: %s needs %s=docker.io/<path> in hack/versions.env\n' \
                    "${file}" "${line}" "${key}" "${pin}" >&2
                status=1
            elif [[ "${!pin}" != docker.io/* ]]; then
                printf '%s:%s: %s uses a registry variable, but %s is %s\n' \
                    "${file}" "${line}" "${key}" "${pin}" "${!pin}" >&2
                status=1
            fi
        done < <(registry_values "${root}/${file}" "${registry}")
    done
}

# Optional arguments: the CI file that sets the image variables, the variable
# that those images start with, and more CI files whose references to it are
# checked the same way.
# Extra registry checks need the image-variable file; do not silently skip them.
if [[ -z "${1:-}" ]]; then
    for argument in "${@:2}"; do
        if [[ -n "${argument}" ]]; then
            echo "IMAGE_VARIABLES_FILE is required when registry checks are requested" >&2
            exit 2
        fi
    done
fi

if (($# > 0)) && [[ -n "$1" ]]; then
    check_image_variables "$1" "${2:-}"

    # check_image_variables already reported a bad registry variable name.
    if [[ "${2:-}" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
        check_registry_references "$2" "$1" "${@:3}"
    fi
fi

exit "${status}"
