#!/bin/sh

# Start Go server in background
/usr/bin/server &

# Wait a moment for server to start
sleep 2

# Start Caddy
exec caddy run --config /etc/caddy/Caddyfile

