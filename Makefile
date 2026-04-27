run:
	go run ./cmd/pulse/main.go

up:
	docker compose up -d

down:
	docker compose down