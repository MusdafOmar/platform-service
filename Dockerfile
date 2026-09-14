FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/platform-service \
    ./cmd/api


FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app/data \
    && chown -R app:app /app

WORKDIR /app

COPY --from=builder /out/platform-service /app/platform-service

USER app

ENV DB_DRIVER=sqlite
ENV DB_DSN=/app/data/platform-service.db

EXPOSE 8080

ENTRYPOINT ["/app/platform-service"]