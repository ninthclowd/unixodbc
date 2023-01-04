all: build

build: .build-api

clean: .clean-api

.build-api:
	DOCKER_BUILDKIT=1 docker build -t go-unixodbc:build --target export-api . --output internal/api

.clean-api:
	rm -f internal/api/cgo_helpers.go internal/api/cgo_helpers.h internal/api/cgo_helpers.c
	rm -f internal/api/const.go internal/api/doc.go internal/api/types.go
	rm -f internal/api/api.go

.snapshot-maria:
	docker compose build mariadb-snapshots
	docker compose run mariadb-snapshots

test:
	docker compose build mariadb-tests
	docker compose run mariadb-tests

snapshots: .snapshot-maria



