api: .clean-api
	DOCKER_BUILDKIT=1 docker build -t go-unixodbc:build --target export-api -f build/Dockerfile . --output internal/api

clean: .clean-api

.clean-api:
	rm -f internal/api/cgo_helpers.go internal/api/cgo_helpers.h internal/api/cgo_helpers.c
	rm -f internal/api/const.go internal/api/doc.go internal/api/types.go
	rm -f internal/api/api.go


trace-mariadb:
	DOCKER_BUILDKIT=1 docker build -t go-unixodbc:mariadb-trace --target export-trace -f test/mariadb/build/Dockerfile . --output .
	go tool trace trace.out

test-mariadb:
	docker build -t go-unixodbc:mariadb-test --target test -f test/mariadb/build/Dockerfile .
	docker run -it --rm  go-unixodbc:mariadb-test

snapshot-mariadb:
	DOCKER_BUILDKIT=1 docker build -t go-unixodbc:mariadb-snapshot --target export-snapshots -f test/mariadb/build/Dockerfile . --output test/mariadb/testdata

test: test-mariadb

snapshots: snapshot-mariadb

mocks:
	go generate ./...






