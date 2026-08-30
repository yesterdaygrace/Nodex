# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder

WORKDIR /app

# deps first for layer cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Go binary (main.go at repo root)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server .

# Build frontend if dist missing (avoids slow npm on every build when host already built it)
RUN if [ ! -f frontend/dist/index.html ] && [ -f frontend/package.json ]; then \
      echo "frontend/dist missing - building..." && \
      apk add --no-cache nodejs npm && \
      cd frontend && npm ci && npm run build || echo "frontend build skipped"; \
    else echo "using existing frontend/dist"; fi
RUN mkdir -p /app/frontend/dist

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app \
    && apk add --no-cache wget ca-certificates

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/frontend/dist /app/frontend/dist

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/server"]
