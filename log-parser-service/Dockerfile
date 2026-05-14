# ----- Build stage -----
FROM golang:1.23-alpine AS build

WORKDIR /src/log-parser-service

RUN apk add --no-cache ca-certificates

COPY go.mod ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /out/log-parser-service ./cmd/log-parser


# ----- Run stage -----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
    && addgroup -S appgroup \
    && adduser -S -H -u 10001 -G appgroup appuser

WORKDIR /app

COPY --from=build /out/log-parser-service ./log-parser-service

RUN chmod 755 /app/log-parser-service \
    && chown -R appuser:appgroup /app

EXPOSE 8080

USER appuser

ENTRYPOINT ["/app/log-parser-service"]
