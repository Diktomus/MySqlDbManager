
MIGRATIONS_CONFIG=config/migrations/dbconfig.yaml

build:
	go build cmd/main.go

migrate-up:
	sql-migrate up -config $(MIGRATIONS_CONFIG)

migrate-down:
	sql-migrate down -config $(MIGRATIONS_CONFIG)
