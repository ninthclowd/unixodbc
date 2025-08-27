

build:
	docker compose down --remove-orphans
	docker compose build

fmt:
	go fmt ./...

test: test-go test-mariadb test-postgres test-mssql

test-go:
	go test ./ ./internal/cache -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out


test-mariadb:
	docker compose run --rm mariadb_test

test-postgres:
	docker compose run --rm postgres_test

test-mssql:
	docker compose run --rm mssql_test