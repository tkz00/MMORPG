db-up:
	docker compose up -d db

db-down:
	docker compose down

run:
	go run -C backend cmd/server/main.go
