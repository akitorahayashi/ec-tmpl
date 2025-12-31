FROM golang:1.22-alpine AS base
WORKDIR /app
RUN apk add --no-cache ca-certificates

FROM base AS development
RUN go install github.com/air-verse/air@v1.52.3
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV EC_TMPL_BIND_IP=0.0.0.0
ENV EC_TMPL_BIND_PORT=8000
EXPOSE 8000
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ec-tmpl ./cmd/ec-tmpl

FROM alpine:3.20 AS production
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/ec-tmpl /app/ec-tmpl
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENV EC_TMPL_BIND_IP=0.0.0.0
ENV EC_TMPL_BIND_PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/entrypoint.sh"]
