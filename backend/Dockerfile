FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM caddy:2.8.0-alpine
COPY --from=builder /app/server /usr/bin/server
COPY Caddyfile /etc/caddy/Caddyfile
COPY start.sh /start.sh
RUN chmod +x /usr/bin/server /start.sh
EXPOSE 8080
CMD ["/start.sh"]
