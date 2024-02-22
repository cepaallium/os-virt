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
FROM docker-cd.unicloud.com/helm/cloud/alpine:3.12.3
WORKDIR /usr/local/
COPY   --from=builder /workspace/os-virt  /usr/local/
COPY   --from=builder /workspace/config  /usr/local/config
ENTRYPOINT ["/usr/local/os-virt"]
