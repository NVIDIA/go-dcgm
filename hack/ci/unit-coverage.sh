#!/usr/bin/env bash

# Generate coverage and JUnit artifacts for the CI-safe unit-test targets.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT_DIR}"

COVERAGE_TARGETS=(
    //cmd/gen-fields:gen-fields_test
    //pkg/dcgm:dcgm_unit_test
    //samples/restApi/handlers:handlers_test
    //samples/restApi:restApi_test
)
TEST_NAMES=(gen-fields dcgm rest-handlers rest-api)
TEST_LOG_PATHS=(
    cmd/gen-fields/gen-fields_test
    pkg/dcgm/dcgm_unit_test
    samples/restApi/handlers/handlers_test
    samples/restApi/restApi_test
)

ARTIFACT_STATUS=0
BAZEL_STATUS=0
COVERAGE_PROFILES=()
TEMPORARY_DIRECTORY=""

# cleanup removes the intermediate coverage profile workspace.
cleanup() {
    [[ -z "${TEMPORARY_DIRECTORY}" ]] || rm -rf -- "${TEMPORARY_DIRECTORY}"
}

# prepare_artifact_directory clears the reports produced by a previous coverage run.
prepare_artifact_directory() {
    rm -rf -- test-reports
    rm -f -- tests.cov coverage.xml
    mkdir -p test-reports
}

# run_bazel_coverage runs every CI-safe unit target while preserving its exit status.
run_bazel_coverage() {
    if bazel coverage \
        --config=ci \
        --cache_test_results=no \
        --keep_going \
        --test_arg=-test.v \
        --@rules_go//go/config:cover_format=go_cover \
        "${COVERAGE_TARGETS[@]}"; then
        BAZEL_STATUS=0
    else
        BAZEL_STATUS=$?
    fi
}

# collect_test_artifacts converts test logs and collects the coverage profiles.
collect_test_artifacts() {
    local index
    local junit
    local profile
    local test_log
    local testlog

    for index in "${!TEST_NAMES[@]}"; do
        testlog="bazel-testlogs/${TEST_LOG_PATHS[index]}"
        test_log="${testlog}/test.log"
        junit="test-reports/unit-${TEST_NAMES[index]}.xml"
        profile="${testlog}/coverage.dat"

        if [[ -s "${test_log}" ]]; then
            if ! go-junit-report -set-exit-code <"${test_log}" >"${junit}"; then
                printf 'failed to convert Bazel test log to JUnit: %s\n' "${test_log}" >&2
                ARTIFACT_STATUS=1
            fi
        else
            printf 'Bazel test log not found: %s\n' "${test_log}" >&2
            ARTIFACT_STATUS=1
        fi

        if [[ -s "${profile}" ]]; then
            COVERAGE_PROFILES+=("${profile}")
        else
            printf 'coverage profile not found: %s\n' "${profile}" >&2
            ARTIFACT_STATUS=1
        fi
    done
}

# exit_on_failures preserves Bazel's failure when both test and artifact work fail.
exit_on_failures() {
    if [[ ${BAZEL_STATUS} -ne 0 ]]; then
        exit "${BAZEL_STATUS}"
    fi
    if [[ ${ARTIFACT_STATUS} -ne 0 ]]; then
        exit "${ARTIFACT_STATUS}"
    fi
}

# merge_coverage_profiles combines overlapping coverprofile blocks into tests.cov.
merge_coverage_profiles() {
    local body
    local mode=""
    local profile
    local profile_mode

    TEMPORARY_DIRECTORY="$(mktemp -d "${TMPDIR:-/tmp}/go-dcgm-unit-coverage.XXXXXX")"
    body="${TEMPORARY_DIRECTORY}/body"

    for profile in "${COVERAGE_PROFILES[@]}"; do
        profile_mode="$(awk 'NR == 1 { print; exit }' "${profile}")"
        if [[ ! "${profile_mode}" =~ ^mode:\ (atomic|count|set)$ ]]; then
            printf 'invalid coverage mode in %s: %s\n' "${profile}" "${profile_mode}" >&2
            exit 1
        fi
        if [[ -n "${mode}" && "${profile_mode}" != "${mode}" ]]; then
            printf 'coverage mode mismatch: %s and %s\n' "${mode}" "${profile_mode}" >&2
            exit 1
        fi
        mode="${profile_mode}"
        awk 'NR > 1 { print }' "${profile}" >>"${body}"
    done

    # Sum atomic/count hits (max for set) the same way a Go coverprofile merger does.
    awk -v mode="${mode}" '
        {
            if (NF < 3) {
                print "malformed coverage line: " $0 > "/dev/stderr"
                exit 1
            }
            key = $1 " " $2
            hits = $3 + 0
            if (!(key in seen)) {
                order[++n] = key
                loc[key] = $1
                stmts[key] = $2
                seen[key] = 1
            }
            if (mode ~ /^mode: set$/) {
                if (hits > count[key]) {
                    count[key] = hits
                }
            } else {
                count[key] += hits
            }
        }
        END {
            print mode
            for (i = 1; i <= n; i++) {
                key = order[i]
                printf "%s %s %.0f\n", loc[key], stmts[key], count[key]
            }
        }
    ' "${body}" >tests.cov
}

# main generates coverage artifacts, then validates and converts the merged profile.
main() {
    trap cleanup EXIT HUP INT TERM
    prepare_artifact_directory
    run_bazel_coverage
    collect_test_artifacts
    exit_on_failures
    merge_coverage_profiles
    go tool cover -func=tests.cov
    gocover-cobertura <tests.cov >coverage.xml
}

main "$@"
