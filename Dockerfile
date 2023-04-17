# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

# caches requirement for runtime/ppnet
COPY runtime/ppnet/go.mod .
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

# caches requirement for app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o bin/backend main.go

CMD ["/app/bin/backend"]

FROM scratch

COPY --from=builder /app/bin/backend /usr/local/bin/backend

CMD ["/usr/local/bin/backend"]