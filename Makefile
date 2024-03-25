

build:
	docker compose down --remove-orphans
	docker compose build

fmt:
	go fmt ./...

test: test-mariadb test-postgres test-mssql

test-mariadb:
	docker compose run --rm mariadb_test
	go tool cover -html test/acceptance/mariadb/coverage.out -o test/acceptance/mariadb/coverage.html

test-postgres:
	docker compose run --rm postgres_test
	go tool cover -html test/acceptance/postgres/coverage.out -o test/acceptance/postgres/coverage.html

test-mssql:
	docker compose run --rm mssql_test
	go tool cover -html test/acceptance/mssql/coverage.out -o test/acceptance/mssql/coverage.html