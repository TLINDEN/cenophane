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
RUN make && strip cenod

FROM alpine:3.17
LABEL maintainer="Uploads Author <info@daemon.de>"

RUN install -o 1001 -g 1001 -d /data

WORKDIR /app
COPY --from=builder /work/cenod /app/cenod

ENV CENOD_LISTEN=:8080
ENV CENOD_STORAGEDIR=/data
ENV CENOD_DBFILE=/data/bbolt.db
ENV CENOD_DEBUG=1

USER 1001:1001
EXPOSE 8080
VOLUME /data
CMD ["/app/cenod"]
