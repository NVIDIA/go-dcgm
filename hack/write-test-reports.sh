#!/usr/bin/env bash

set -euo pipefail

coverage_profiles=()
coverage_run=false

# result_xmls decodes the JUnit paths recorded by Bazel for one invocation.
result_xmls() {
  local event_log="$1"
  # File URIs encode bytes such as spaces as %XX.
  grep -o '"uri":"file://[^"]*/testlogs/[^"]*/test\.xml"' "${event_log}" \
    | sed -e 's|^"uri":"file://||' -e 's|"$||' \
    | awk '
      {
        out = ""
        digits = "0123456789abcdef"
        while (match($0, /%[0-9A-Fa-f][0-9A-Fa-f]/)) {
          hex = tolower(substr($0, RSTART + 1, 2))
          high = index(digits, substr(hex, 1, 1)) - 1
          low = index(digits, substr(hex, 2, 1)) - 1
          out = out substr($0, 1, RSTART - 1) sprintf("%c", high * 16 + low)
          $0 = substr($0, RSTART + 3)
        }
        print out $0
      }' | sort -u || true
}

# collect_test_results copies this run's JUnit files and records its coverage profiles.
collect_test_results() {
  local log coverage_log xmls xml dir target
  for log in "$@"; do
    coverage_log=false
    if grep -q '"command":"coverage"' "${log}"; then
      coverage_log=true
      coverage_run=true
    fi
    xmls="$(result_xmls "${log}")"
    while IFS= read -r xml; do
      if [[ -z "${xml}" ]]; then
        continue
      fi
      dir="$(dirname "${xml}")"
      # The path below testlogs keeps the targets' names apart.
      target="${dir#*/testlogs/}"
      cp "${xml}" "test-reports/$(printf '%s' "${target}" | tr / _).xml"
      if [[ "${coverage_log}" == true ]]; then
        if [[ -s "${dir}/coverage.dat" ]]; then
          coverage_profiles+=("${dir}/coverage.dat")
        else
          echo "warning: ${target}: the coverage run wrote no coverage.dat" >&2
        fi
      fi
    done <<< "${xmls}"
  done
}

# write_coverage merges profiles from executed tests, then writes the CI coverage report.
write_coverage() {
  # Each profile also lists the packages the test only imports, so a
  # block appears in many profiles: sum its counts, or keep any hit in
  # set mode.
  awk '
    FNR == 1 {
      if ($0 !~ /^mode: (set|count|atomic)$/) {
        print "error: " FILENAME ": invalid coverage mode: " $0 > "/dev/stderr"
        bad = 1
      }
      if (mode == "") mode = $0
      else if ($0 != mode) { print "coverage mode mismatch: " mode " and " $0 > "/dev/stderr"; bad = 1 }
      next
    }
    {
      key = $1 " " $2
      if (!(key in hits)) order[++n] = key
      hits[key] = (mode == "mode: set") ? (hits[key] || $3 > 0) : hits[key] + $3
    }
    END {
      if (bad) exit 1
      print mode
      for (i = 1; i <= n; i++) print order[i], hits[order[i]]
    }
  ' "${coverage_profiles[@]}" > tests.cov
  go tool cover -func=tests.cov | tail -n 1
  gocover-cobertura < tests.cov > coverage.xml
}

# main reads the selected Bazel event folder and replaces only generated report files.
main() {
  if [[ "$#" -ne 1 ]]; then
    echo "usage: $0 <bazel-event-directory>" >&2
    exit 2
  fi

  local events_dir="$1"
  local -a logs
  local root
  root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
  cd "${root}"

  rm -rf test-reports tests.cov coverage.xml
  mkdir -p test-reports

  shopt -s nullglob
  logs=("${events_dir}"/*.json)
  if [[ "${#logs[@]}" -eq 0 ]]; then
    echo "error: no Bazel event log in ${events_dir}; run a test task first" >&2
    exit 1
  fi

  collect_test_results "${logs[@]}"

  if [[ -z "$(ls -A test-reports)" ]]; then
    echo "error: the Bazel event logs list no test results" >&2
    exit 1
  fi
  if [[ "${coverage_run}" == false ]]; then
    exit 0
  fi
  if [[ "${#coverage_profiles[@]}" -eq 0 ]]; then
    echo "error: the coverage run wrote no coverage profiles" >&2
    exit 1
  fi

  write_coverage
}

main "$@"
