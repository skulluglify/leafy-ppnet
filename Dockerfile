# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

ENV CGO_ENABLED 0
ENV DOCKERIZED 1

COPY . .
COPY .env.docker /
RUN go build -o bin/backend main.go

CMD ["/app/bin/backend"]

FROM scratch

ENV DOCKERIZED 1

COPY --from=builder /app/bin/backend /usr/local/bin/backend

CMD ["/usr/local/bin/backend"]