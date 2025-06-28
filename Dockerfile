# Build the manager binary
FROM golang:1.23 AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager cmd/main.go

# Use minimal base image
FROM alpine:3.18
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
