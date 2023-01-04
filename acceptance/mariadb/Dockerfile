FROM golang:1.19-bullseye as base
RUN apt update
RUN apt-get install -y unixodbc-dev unixodbc

# install mariadb odbc connector
WORKDIR /usr/src/odbc
ADD https://dlm.mariadb.com/2454035/Connectors/odbc/connector-odbc-3.1.17/mariadb-connector-odbc-3.1.17-debian-bullseye-amd64.tar.gz connector.tar.gz
RUN tar -zxvf connector.tar.gz
RUN mv mariadb-connector* connector
WORKDIR /usr/src/odbc/connector
RUN install lib/mariadb/libmaodbc.so /usr/local/lib/
RUN install lib/mariadb/libmariadb.so.3 /usr/local/lib/
RUN install -d /usr/local/lib/mariadb/
RUN install -d /usr/local/lib/mariadb/plugin/
RUN install lib/mariadb/plugin/* /usr/local/lib/mariadb/plugin/
RUN ldconfig
WORKDIR /usr/src/odbc
RUN echo "[MariaDB ODBC 3.0 Driver]\n" \
         "Description = MariaDB Connector/ODBC v.3.0\n" \
         "Driver = /usr/local/lib/libmaodbc.so\n" > odbcinst.ini
RUN odbcinst -i -d -f odbcinst.ini
RUN echo "[MariaDB]\n" \
         "Description=MariaDB server\n" \
         "Driver=MariaDB ODBC 3.0 Driver\n" \
         "SERVER=mariadb\n" \
         "USER=test\n" \
         "PASSWORD=test\n" > odbc.ini
RUN odbcinst -i -s -l -f odbc.ini

WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG UPDATE_GOLDEN
CMD ["go","test","./acceptance/mariadb/..."]
