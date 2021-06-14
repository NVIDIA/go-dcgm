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

GOLANG_VERSION := 1.14.2

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

test-main:
	go test ./tests

check-format:
	test $$(gofmt -l . | tee /dev/stderr | wc -l) -eq 0

