FROM golang:1.20-bookworm as base
RUN apt update
RUN apt-get install -y libssl-dev unixodbc-dev unixodbc odbcinst odbc-postgresql

WORKDIR /usr/src/unixodbc


