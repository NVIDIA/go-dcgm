# syntax=docker/dockerfile:1.19

# CUDA_VERSION and DISTRO_FLAVOR are used to select a docker image tag from the upstream
# docker registry for nvidia/cuda.   The variation of DISTRO_FLAVOR and CUDA_VERSION must
# point to an image that exists, see here for list: https://hub.docker.com/r/nvidia/cuda/tags

# CUDA_VERSION
ARG CUDA_VERSION=13.3.1
# cuda image supports these images rockylinux9, rockylinux8, ubi9, ubi8, ubuntu26.04, ubuntu24.04, ubuntu22.04
# Note: Testing has only been done with the ubuntu variants.
ARG DISTRO_FLAVOR=ubuntu26.04

# Use build arguments to select our base image or just stick with the defaults above.
FROM nvidia/cuda:$CUDA_VERSION-base-$DISTRO_FLAVOR AS base
ARG DCGM_VERSION=4.6.1-1
ARG DCGM_LOCAL_REPO_DEB_URL=""
ARG GO_VERSION=1.26.7
ENV DEBIAN_FRONTEND=noninteractive

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup our apt environment and install the necessary keyrings and repositories to install dcgm.  Note that this strategy doesn't
# support dcgm 3.x.
# Keep DCGM packages on the same version to avoid apt selecting a newer core package.
# hadolint ignore=DL3015,DL3008
RUN --mount=type=secret,id=gitlab_job_token,required=false \
    apt-get update && apt-get install -y --no-install-recommends \
    gnupg2 curl ca-certificates git && \
    if [[ -n "${DCGM_LOCAL_REPO_DEB_URL}" ]]; then \
      token_args=(); \
      if [[ -f /run/secrets/gitlab_job_token ]]; then \
        token_args=(-H "JOB-TOKEN: $(cat /run/secrets/gitlab_job_token)"); \
      fi; \
      curl -fsSL "${token_args[@]}" -o /tmp/dcgm-local-repo.deb "${DCGM_LOCAL_REPO_DEB_URL}"; \
      apt-get install -y --no-install-recommends /tmp/dcgm-local-repo.deb; \
      rm -f /tmp/dcgm-local-repo.deb; \
      apt-get update; \
    fi && \
    go_tarball="go${GO_VERSION}.linux-amd64.tar.gz" && \
    curl -fsSLo "/tmp/${go_tarball}" "https://go.dev/dl/${go_tarball}" && \
    go_sha256="$(curl -fsSL "https://dl.google.com/go/${go_tarball}.sha256")" && \
    printf '%s  %s\n' "${go_sha256}" "/tmp/${go_tarball}" | sha256sum -c - && \
    tar -C /usr/local -xzf "/tmp/${go_tarball}" && \
    rm "/tmp/${go_tarball}" && \
    apt-get install -y \
      datacenter-gpu-manager-4-core=1:${DCGM_VERSION} \
      datacenter-gpu-manager-4-dev=1:${DCGM_VERSION} && \
    apt-get purge --autoremove -y curl && \
    rm -rf /var/lib/apt/lists/*

ENV PATH=$PATH:/usr/local/go/bin

# build go-dcgm and samples inside docker environment
FROM base AS samples
# hadolint ignore=DL3008,DL3015
RUN apt-get update && apt-get install -y build-essential nvidia-utils-610 && rm -rf /var/lib/apt/lists/*
COPY . /src
WORKDIR /src
RUN make binary && \
    cp ./samples/restApi/restApi \
      ./samples/processInfo/processInfo \
      ./samples/diag/diag \
      ./samples/hostengineStatus/hostengineStatus \
      ./samples/dmon/dmon \
      ./samples/health/health \
      ./samples/topology/topology \
      ./samples/deviceInfo/deviceInfo \
      ./samples/policy/policy \
    /usr/local/go/bin/
WORKDIR /
