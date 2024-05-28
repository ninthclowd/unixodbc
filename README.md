# Go unixodbc driver

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Build](https://github.com/ninthclowd/unixodbc/actions/workflows/build.yml/badge.svg?branch=main)
![Acceptance Tests](https://github.com/ninthclowd/unixodbc/actions/workflows/acceptance.yml/badge.svg?branch=main)
![CodeQL](https://github.com/ninthclowd/unixodbc/actions/workflows/codeql.yml/badge.svg?branch=main)

This topic provides instructions for installing, running, and modifying the go unixodbc driver for connecting to ODBC databases
through [unixodbc](https://www.unixodbc.org/). The driver supports Go's [database/sql](https://golang.org/pkg/database/sql/) package.

# Prerequisites

The following software packages are required to use the go unixodbc driver.

## Go

The latest driver requires the [Go language](https://golang.org/) 1.20 or higher.

## Operating System

This driver was primarily developed with support of [Debian 12](https://www.debian.org/releases/bookworm/), however other
linux distributions may work correctly providing that [unixodbc](#unixodbc) is installed.

## [unixodbc](https://www.unixodbc.org/)

[unixodbc](https://www.unixodbc.org/) 2.3.11 or greater must on the system your application is running on. For debian, the following packages must be installed in order for this
driver to connect:

- libssl-dev
- unixodbc-common
- unixodbc-dev
- unixodbc

## Supported Databases and Database Drivers

The driver has been developed and tested for the following databases with the corresponding driver, but may work with other proprietary databases:

| Database                                | Tested Version | Database Driver                                                   |
| --------------------------------------- | -------------- | ----------------------------------------------------------------- |
| [postgres](https://www.postgresql.org/) | 16             | [odbc-postgresql](https://odbc.postgresql.org/)                   |
| Microsoft SQL Server                    | 2022           | [tdsodbc](https://www.freetds.org/)                               |
| [mariadb](https://mariadb.org/)         | 10.9           | [odbc-mariadb](https://mariadb.com/kb/en/mariadb-connector-odbc/) |

# Configuration

## Basic Connection

Typical connection to the database can be established by importing the go unixodbc driver and opening the connection with
[sql.Open](https://pkg.go.dev/database/sql#Open).
The connection string to use will be specific to the [database and database driver](#supported-databases-and-database-drivers)
that you are using:

```go
package main

import (
    _ "github.com/ninthclowd/unixodbc"
    "database/sql"
)

func main(){
    db = sql.Open("unixodbc", "DSN=EXAMPLE")
}
```

## Prepared Statement Caching

The driver supports prepared statement caching on each connection using a LRU algorithm by connecting to the database
with [sql.OpenDB](https://pkg.go.dev/database/sql#OpenDB) and supplying an
[unixodbc.Connector](https://pkg.go.dev/github.com/ninthclowd/unixodbc#Connector):

```go
package main

import (
    "github.com/ninthclowd/unixodbc"
    "database/sql"
)

func main(){
    db = sql.OpenDB(&unixodbc.Connector{
        ConnectionString:   unixodbc.StaticConnStr("DSN=EXAMPLE"),
        StatementCacheSize: 5, // number of prepared statements to cache per connection
    })
}
```

## Dynamic Connection Strings

The driver supports dynamic connection strings for use in databases that require token based authentication. This can
be accomplished by implementing [unixodbc.ConnectionStringFactory](https://pkg.go.dev/github.com/ninthclowd/unixodbc#ConnectionStringFactory)
and connecting with [sql.OpenDB](https://pkg.go.dev/database/sql#OpenDB):

```go
package main

import (
	"github.com/ninthclowd/unixodbc"
	"database/sql"
)

type GetToken struct {
}

// ConnectionString implements unixodbc.ConnectionStringFactory
func (d *GetToken) ConnectionString() (connStr string, err error) {
	var token string
	// insert code to pull token for each new connection
	// ...

	// create a dynamic connection string and return it
	connStr = "DSN=EXAMPLE;UID=User;PWD=" + token
	return
}

func main() {
	db = sql.OpenDB(&unixodbc.Connector{
		ConnectionString: &GetToken{},
	})
}
```

# Development

To develop this code base you will need the following components installed on your system:

- standard C headers
- the [unixodbc](#unixodbc) C headers
- [mockgen](https://github.com/golang/mock)
- [c-for-go](https://github.com/xlab/c-for-go)

API wrappers and mocks for unit testing are generated through `go generate ./...`

An example for setting up an Ubuntu WSL environment:
```bash
apt install gcc-multilib libssl-dev unixodbc-dev
go install github.com/golang/mock/mockgen@v1.6.0
go install https://github.com/xlab/c-for-go@v1.2.0
```

## Acceptance Testing

At this time, most CGO driver functionality in the [internal/odbc package](internal/odbc) is validated through the acceptance tests in the [test/acceptance](test/acceptance) folder using the corresponding [workflow](.github/workflows/acceptance.yml) and [docker-compose.yml](docker-compose.yml).
To run these tests, ensure [Docker](https://www.docker.com/) is installed and run `make test`.

## Submitting Pull Requests

You may use your preferred editor to edit the driver code. If you are adding support for a new database, please include
new [Acceptance Tests](#acceptance-testing) for the database if possible. Please `make fmt` to format your code
and `make test` to validate your code before submitting.
