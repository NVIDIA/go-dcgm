# Copyright (c) 2020, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOLANG_VERSION := 1.23.6
GOLANGCILINT_TIMEOUT ?= 10m

.PHONY: all binary install check-format
all: binary test-main check-format

binary:
	go build ./pkg/dcgm
	cd samples/deviceInfo; go build
	cd samples/dmon; go build
	cd samples/health; go build
	cd samples/hostengineStatus; go build
	cd samples/policy; go build
	cd samples/processInfo; go build
	cd samples/restApi; go build
	cd samples/topology; go build
	cd samples/diag; go build

test-main:
	go test -race ./tests

check-format:
	test $$(gofumpt -l -w . | tee /dev/stderr | wc -l) -eq 0

clean:
	rm -f samples/deviceInfo/deviceInfo
	rm -f samples/dmon/dmon
	rm -f samples/health/health
	rm -f samples/hostengineStatus/hostengineStatus
	rm -f samples/policy/policy
	rm -f samples/processInfo/processInfo
	rm -f samples/restApi/restApi
	rm -f samples/topology/topology

lint:
	golangci-lint run ./... --timeout $(GOLANGCILINT_TIMEOUT)  --new-from-rev=HEAD~1 --fix

lint-full:
	golangci-lint run ./... --timeout $(GOLANGCILINT_TIMEOUT) --fix
