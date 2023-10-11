include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATIONS_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DB) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.15.1

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v