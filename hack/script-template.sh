#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: Apache-2.0
#
# Starting point for a new script. Copy it, replace this paragraph with what
# the script does, and replace the example flags, check_flags and main.
set -euo pipefail

arg0="$(basename "${0}")"
blank="${arg0//?/ }"

# Flag defaults. Replace the examples with the script's own flags.
example_text=""
example_count=1

# usage_info prints the usage lines. Keep them in step with flags.
usage_info() {
  echo "Usage: ${arg0} --text TEXT [--count COUNT]"
  echo "       ${blank} [--help]"
}

# usage prints an optional message and the usage lines to standard error,
# then exits 2. Call it for any mistake on the command line.
usage() {
  exec 1>&2
  if [[ "${#}" -gt 0 ]]; then
    echo "${arg0}: ${*}"
  fi
  usage_info
  exit 2
}

# help prints the usage lines and what each flag does, then exits 0.
help() {
  usage_info
  echo ""
  echo "Options:"
  echo "  -t, --text TEXT     text to print"
  echo "  -c, --count COUNT   how many times to print it (default: 1)"
  echo "  -h, --help          show this help message"
  exit 0
}

# error prints a message to standard error and exits 1. Call it when the
# work fails.
error() {
  echo "${arg0}: ${*}" >&2
  exit 1
}

# flags reads the command line into the flag variables.
flags() {
  while [[ "${#}" -gt 0 ]]; do
    case "${1}" in
      -t | --text)
        [[ "${#}" -gt 1 ]] || usage "${1} needs a value"
        example_text="${2}"
        shift 2
        ;;
      -c | --count)
        [[ "${#}" -gt 1 ]] || usage "${1} needs a value"
        example_count="${2}"
        shift 2
        ;;
      -h | --help) help ;;
      --)
        shift
        break
        ;;
      -*) usage "unknown option ${1}" ;;
      *) break ;;
    esac
  done

  # This script takes no other arguments. Read them here if yours does.
  [[ "${#}" -eq 0 ]] || usage "unexpected argument ${1}"
}

# check_flags checks the flag values once all of them are read.
check_flags() {
  [[ -n "${example_text}" ]] || usage "--text is required"
  # Nine digits at most, so the number cannot overflow shell arithmetic.
  [[ "${example_count}" =~ ^[0-9]{1,9}$ ]] || usage "--count needs a whole number of at most 9 digits"
  # Read the number as base 10, so a leading zero does not make it octal.
  example_count=$((10#${example_count}))
}

# main does the work. Replace the example.
main() {
  local i

  for ((i = 0; i < example_count; i++)); do
    printf '%s\n' "${example_text}"
  done
}

flags "${@}"
check_flags
main
