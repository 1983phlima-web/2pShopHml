# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

COPY go.mod ./
RUN go mod tidy && go mod download

COPY . .

ARG APP=api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.version=$(git describe --tags --always 2>/dev || echo 'dev') -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse --short HEAD 2>/dev || echo 'unknown')" -o /bin/app ./apps/${APP}

# Runtime stage
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /bin/app /app

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app"]
