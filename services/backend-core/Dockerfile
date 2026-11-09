# ==========================================
# STAGE 1: BUILD
# ==========================================
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .

# Fetch dependencies & compile static binary
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server cmd/api/main.go

# ==========================================
# STAGE 2: PRODUCTION RUNTIME
# ==========================================
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl
ENV TZ=Asia/Jakarta

COPY --from=builder /app/server /app/server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/server"]
