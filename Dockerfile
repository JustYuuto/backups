FROM golang:latest as api-builder

WORKDIR /build

RUN cd api && \
    go mod tidy && \
    go build -o api .

FROM oven/bun:latest as web-builder

WORKDIR /build

RUN cd web && \
    bun install --frozen-lockfile && \
    bun run build

FROM debian:stable-slim

WORKDIR /app

# Install CA certificates
RUN apt update && apt install -y ca-certificates

COPY --from=api-builder /build/api /app/api
COPY --from=web-builder /build/dist /app/web

ENV GIN_MODE=release

CMD ["/app/api"]