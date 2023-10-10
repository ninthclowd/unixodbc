FROM golang:1.21-bookworm as base
RUN apt update
RUN apt-get install -y libssl-dev unixodbc-dev unixodbc odbcinst odbc-mariadb

WORKDIR /usr/src/odbc
COPY test/acceptance/mariadb/odbc.ini .
RUN odbcinst -i -s -l -f odbc.ini

WORKDIR /usr/src/unixodbc



