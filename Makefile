include .envrc

DB_DSN=$(INVENTORY_DB_DSN)

.PHONY: help
help:
	@echo 'Available commands:'
	@echo '  make db/up         - Apply all pending migrations'
	@echo '  make db/down n=1   - Rollback last N migrations'
	@echo '  make db/version    - Show current migration version'
	@echo '  make db/force v=5  - Force version (fix dirty state)'
	@echo '  make db/psql       - Connect to database'
	@echo '  make db/reset      - Drop and recreate database'
	@echo '  make run           - Run the Go application'

.PHONY: db/up
db/up:
	migrate -path ./internal/db/migrations -database "$(DB_DSN)" up

.PHONY: db/down
db/down:
	migrate -path ./internal/db/migrations -database "$(DB_DSN)" down $(n)

.PHONY: db/version
db/version:
	migrate -path ./internal/db/migrations -database "$(DB_DSN)" version

.PHONY: db/force
db/force:
	migrate -path ./internal/db/migrations -database "$(DB_DSN)" force $(v)

.PHONY: db/psql
db/psql:
	psql "$(DB_DSN)"

.PHONY: db/reset
db/reset:
	@echo 'Dropping and recreating database...'
	-psql -U postgres -c "DROP DATABASE IF EXISTS inventory_db;"
	psql -U postgres -c "CREATE DATABASE inventory_db OWNER vinfa;"
	@echo 'Database reset. Run "make db/up" to apply migrations.'

.PHONY: run
run:
	go run ./cmd/api
