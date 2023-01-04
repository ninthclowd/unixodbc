FROM golang:bullseye as base
RUN apt update
RUN apt install -y unixodbc-dev
RUN go install github.com/xlab/c-for-go@latest
WORKDIR /usr/src/app
COPY api.yml .
RUN c-for-go api.yml

FROM scratch AS export-api
COPY --from=base /usr/src/app/api/* .

