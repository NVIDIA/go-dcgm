# syntax=docker/dockerfile:1.28.0@sha256:bb22d9815c728170f72750f4e5b0d672e06176142e1d602c7e66c050100b7e5b

FROM nvcr.io/nvidia/cuda:13.4.3-base-ubuntu26.04@sha256:e2eb43e25d404e2c48716f80afdb8989a84c4da588f5509330322842084d97b0

ARG DCGM_VERSION=4.7.0-2

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

COPY hack/versions.env /opt/go-dcgm/hack/versions.env

RUN usermod --login developer --home /home/developer --move-home ubuntu \
    && groupmod --new-name developer ubuntu

# Keep all DCGM packages at the selected version.
# hadolint ignore=DL3008,DL3015
RUN set -eux; \
    source /opt/go-dcgm/hack/versions.env; \
    apt-get update; \
    printf '%s\n' \
        'Package: datacenter-gpu-manager-4* libdcgm dcgmi nv-hostengine' \
        "Pin: version 1:${DCGM_VERSION}" \
        'Pin-Priority: 1000' \
        > /etc/apt/preferences.d/dcgm; \
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        git \
        pkg-config \
        shellcheck \
        "datacenter-gpu-manager-4-core=1:${DCGM_VERSION}" \
        "datacenter-gpu-manager-4-cuda13=1:${DCGM_VERSION}" \
        "datacenter-gpu-manager-4-dev=1:${DCGM_VERSION}"; \
    rm /etc/apt/preferences.d/dcgm; \
    case "$(dpkg --print-architecture)" in \
        amd64) go_archive="go${GO_VERSION}.linux-amd64.tar.gz"; go_sha256="${GO_LINUX_AMD64_SHA256}" ;; \
        arm64) go_archive="go${GO_VERSION}.linux-arm64.tar.gz"; go_sha256="${GO_LINUX_ARM64_SHA256}" ;; \
        *) echo "unsupported architecture: $(dpkg --print-architecture)" >&2; exit 1 ;; \
    esac; \
    curl --fail --silent --show-error --location --output "/tmp/${go_archive}" "https://dl.google.com/go/${go_archive}"; \
    printf '%s  %s\n' "${go_sha256}" "/tmp/${go_archive}" | sha256sum --check --status; \
    tar --directory /usr/local --extract --gzip --file "/tmp/${go_archive}"; \
    rm "/tmp/${go_archive}"; \
    export GOBIN=/usr/local/bin GOMODCACHE=/tmp/go-mod-cache GOCACHE=/tmp/go-build-cache; \
    /usr/local/go/bin/go install "github.com/go-task/task/v3/cmd/task@v${TASK_VERSION}"; \
    /usr/local/go/bin/go install "github.com/bazelbuild/bazelisk@v${BAZELISK_VERSION}"; \
    /usr/local/go/bin/go install "mvdan.cc/gofumpt@v${GOFUMPT_VERSION}"; \
    /usr/local/go/bin/go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}"; \
    /usr/local/go/bin/go install "github.com/boumenot/gocover-cobertura@v${GOCOVER_COBERTURA_VERSION}"; \
    /usr/local/go/bin/go install "github.com/jstemmer/go-junit-report/v2@v${GO_JUNIT_REPORT_VERSION}"; \
    ln --symbolic /usr/local/bin/bazelisk /usr/local/bin/bazel; \
    mkdir /opt/bazelisk; \
    BAZELISK_HOME=/opt/bazelisk USE_BAZEL_VERSION="${BAZEL_VERSION}" bazelisk version; \
    chmod --recursive a+rX /opt/bazelisk; \
    rm --recursive --force /tmp/go-build-cache /tmp/go-mod-cache /var/lib/apt/lists/*

ENV BAZELISK_HOME=/opt/bazelisk \
    GOTOOLCHAIN=local \
    PATH=/usr/local/go/bin:/usr/local/bin:${PATH}

COPY .bazelrc .bazelversion BUILD.bazel MODULE.bazel MODULE.bazel.lock go.mod go.sum /opt/go-dcgm/

WORKDIR /opt/go-dcgm

RUN bazel mod deps && bazel fetch //:gazelle

WORKDIR /
