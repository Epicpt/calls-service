.PHONY:migrate-up migrate-down migrate-status

MIGRATIONS_DIR = ./migrations

DB_DSN = postgres://postgres:1@localhost:5432/callsdb?sslmode=disable

migrate-auth-up:
	migrate -path $(MIGRATIONS_DIR)/auth -database $(DB_DSN) up

migrate-auth-down:
	migrate -path $(MIGRATIONS_DIR)/auth -database $(DB_DSN) down

migrate-calls-up:
	migrate -path $(MIGRATIONS_DIR)/calls -database $(DB_DSN) up

migrate-calls-down:
	migrate -path $(MIGRATIONS_DIR)/calls -database $(DB_DSN) down

generate-mocks:
	@echo "Start generating mocks..."
	@mockery --config .mockery.yaml

MOCKS_DIR ?= ./mocks
clean:
	@echo "Start clean mocks..."
	@rm -rf $(MOCKS_DIR)/