

build:
	docker compose down --remove-orphans
	docker compose build

test: test-mariadb

test-mariadb:
	docker compose run --rm mariadb_test
	go tool cover -html test/acceptance/mariadb/coverage.out -o test/acceptance/mariadb/coverage.html

