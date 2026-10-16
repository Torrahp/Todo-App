-include .env
export

TEST_DB_URL = postgres://postgres:$(POSTGRES_PASSWORD)@db:5432/todoapp_test?sslmode=disable

.PHONY: up down reset logs migrate migrate-test test test-integration

up:
	docker compose up --build -d

down:
	docker compose down

reset:
	docker compose down -v

logs:
	docker compose logs -f backend

migrate:
	docker compose run --rm migrate

migrate-test:
	docker compose run --rm migrate -path=/migrations -database "$(TEST_DB_URL)" up

test:
	cd backend && go test ./...
	cd frontend && npx vitest run

test-integration: migrate-test
	cd backend && DB_HOST=localhost DB_PORT=5433 DB_USER=postgres DB_PASSWORD=$(POSTGRES_PASSWORD) DB_NAME=todoapp_test \
		go test -tags=integration -count=1 ./internal/repository/...
