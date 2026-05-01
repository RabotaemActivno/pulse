DB_MIGRATE_URL = postgres://login:pass@localhost:5432/postgres?sslmode=disable
MIGRATE_PATH = ./migration/postgres

run:
	go run ./cmd/pulse/main.go

up:
	docker compose up -d --build

down:
	docker compose down

migrate-create:
	@read -p "Name:" name; \
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $$name

migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all