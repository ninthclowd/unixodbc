FROM golang:1.21-bookworm as base
RUN apt update
RUN apt-get install -y libssl-dev unixodbc-dev unixodbc odbcinst odbc-postgresql

WORKDIR /usr/src/odbc
COPY test/acceptance/postgres/odbc.ini .
RUN odbcinst -i -s -l -f odbc.ini

WORKDIR /usr/src/unixodbc


