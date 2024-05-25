# syntax = docker/dockerfile:1-experimental

ARG GO_VERSION=1.22.0
ARG BB_VERSION=1.36.1
FROM golang:${GO_VERSION} as builder

WORKDIR /app

COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build make build

FROM busybox:${BB_VERSION}

RUN mkdir -p /lib/time
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /lib/time
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /app/dot-proxy /app/dot-proxy

ENV ZONEINFO=/lib/time/zoneinfo.zip

WORKDIR /app
RUN chmod +x /app/dot-proxy

ENTRYPOINT [ "/app/dot-proxy" ] 