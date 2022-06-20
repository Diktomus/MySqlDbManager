APP=mysql-dbmanager
MIGRATIONS_CONFIG=config/migrations/dbconfig.yaml

run-docker:
	docker-compose up

build-docker:
	docker-compose build

run:
	bin/$(APP)

build:
	go build -o bin/$(APP) cmd/main.go

migrate-up:
	sql-migrate up -config $(MIGRATIONS_CONFIG)

migrate-down:
	sql-migrate down -config $(MIGRATIONS_CONFIG)
