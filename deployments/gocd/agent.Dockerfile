FROM golang:1.27-alpine AS go-toolchain

FROM gocd/gocd-agent-wolfi:v26.1.0

USER root

COPY --from=go-toolchain /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}"

RUN apk add --no-cache \
    docker-cli \
    make

USER go