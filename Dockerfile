# CUDA_VERSION and DISTRO_FLAVOR are used to select a docker image tag from the upstream
# docker registry for nvidia/cuda.   The variation of DISTRO_FLAVOR and CUDA_VERSION must
# point to an image that exists, see here for list: https://hub.docker.com/r/nvidia/cuda/tags

# CUDA_VERSION
ARG CUDA_VERSION=12.5.1
# cuda image supports these images rockylinux9, rockylinux8, ubi9, ubi8, ubuntu24.04, ubuntu22.04, ubuntu20.04
# Note: Testing has only been done with the ubuntu variants.
ARG DISTRO_FLAVOR=ubuntu24.04

# Use build arguments to select our base image or just stick with the defaults above.
FROM nvidia/cuda:$CUDA_VERSION-base-$DISTRO_FLAVOR AS base
ARG DCGM_VERSION=4.2.3-2
ARG GO_VERSION=1.24.4
ENV DEBIAN_FRONTEND=noninteractive

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup our apt environment and install the necessary keyrings and repositories to install dcgm.  Note that this strategy doesn't
# support dcgm 3.x.
# We want recommended packages for dcgm and we dont want to enforce version pinning...yet
# hadolint ignore=DL3015,DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    gnupg2 curl ca-certificates && \
    curl -fsSL https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2004/x86_64/cuda-keyring_1.1-1_all.deb | apt-get install -y --no-install-recommends && \
    curl -fsSL https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu2004/x86_64/nvidia-machine-learning-repo-ubuntu2004_1.0.0-1_amd64.deb | apt-get install -y --no-install-recommends && \
    curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz && \
    apt-get purge --autoremove -y curl && \
    apt-get install -y datacenter-gpu-manager-4-dev=1:${DCGM_VERSION} && \
    rm -rf /var/lib/apt/lists/*

ENV PATH=$PATH:/usr/local/go/bin

# build go-dcgm and samples inside docker environment
FROM base AS samples
# hadolint ignore=DL3008,DL3015
RUN apt-get update && apt-get install -y build-essential nvidia-utils-555 && rm -rf /var/lib/apt/lists/*
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
