# syntax=docker/dockerfile:1 -*-python-*-
FROM golang:1.20-alpine as builder

RUN apk update
RUN apk upgrade
RUN apk add --no-cache git make bash binutils

RUN git --version

WORKDIR /work

COPY go.mod .
RUN go mod download
COPY . .
RUN make && strip upd

FROM alpine:3.17
LABEL maintainer="Uploads Author <info@daemon.de>"

RUN install -o 1001 -g 1001 -d /data

WORKDIR /app
COPY --from=builder /work/upd /app/upd

ENV UPD_LISTEN=:8080
ENV UPD_STORAGEDIR=/data
ENV UPD_DBFILE=/data/bbolt.db
ENV UPD_DEBUG=1

USER 1001:1001
EXPOSE 8080
VOLUME /data
CMD ["/app/upd"]
