# syntax=docker/dockerfile:1.19

FROM nvcr.io/nvidia/cuda:13.4.1-base-ubuntu26.04@sha256:9b091bbc3e1481238346e90d93f63684af509c08ab32218f98793141fc540fd9

ARG DCGM_VERSION=4.7.0-1

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

COPY hack/versions.env /opt/go-dcgm/hack/versions.env

RUN usermod --login developer --home /home/developer --move-home ubuntu \
    && groupmod --new-name developer ubuntu

# hadolint ignore=DL3008,DL3015
RUN set -eux; \
    source /opt/go-dcgm/hack/versions.env; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        git \
        pkg-config \
        shellcheck \
        "datacenter-gpu-manager-4-core=1:${DCGM_VERSION}" \
        "datacenter-gpu-manager-4-dev=1:${DCGM_VERSION}"; \
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
