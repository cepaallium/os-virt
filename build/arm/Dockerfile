#Copyright 2021 KubeCube Authors
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

# Build the manager binary
FROM docker-cd.unicloud.com/library/golang:1.17 as builder

WORKDIR /workspace

# Copy the go source
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY config/ config/

# Build
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on GOPROXY=http://10.210.20.31/repository/goproxy/ go build -a -o os-virt cmd/os-virt.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM docker-cd.unicloud.com/arm64/helm/cloud/alpine:3.12.3
WORKDIR /usr/local/
COPY   --from=builder /workspace/os-virt  /usr/local/
COPY   --from=builder /workspace/config  /usr/local/config
ENTRYPOINT ["/usr/local/os-virt"]
