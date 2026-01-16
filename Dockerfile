# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Transfer source code
COPY src ./src
COPY *.go ./

# Build
RUN CGO_ENABLED=0 go build -o /dist/dispatcher

# Test
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine/curl AS build-release-stage

WORKDIR /app

COPY --from=builder /dist /app

ENTRYPOINT ["/app/dispatcher"]
