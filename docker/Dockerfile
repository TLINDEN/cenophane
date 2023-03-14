# syntax=docker/dockerfile:1 -*-python-*-
FROM golang:1.18-alpine
ARG VERSION
LABEL upd_version=0.1

ENV HOME="/" \
    OS_ARCH="amd64" \
    OS_FLAVOUR="alpine" \
    OS_NAME="linux"

RUN echo $VERSION
RUN install -o 1001 -g 1001 -d /data

WORKDIR /build
RUN mkdir -p api cfg cmd
COPY upd/api/* ./api/
COPY upd/cfg/* ./cfg/
COPY upd/cmd/* ./cmd/
COPY upd/*.* ./
RUN go mod tidy
#RUN echo $VERSION
# doesn't work, see
# https://blog.alexellis.io/inject-build-time-vars-golang/
RUN go build -ldflags "-X upd/cfg.VERSION=$VERSION"
RUN echo $VERSION > version
RUN ls -l


WORKDIR /app
RUN cp /build/upd /app/upd
#RUN rm -rf /build

ENV UPD_LISTEN=:8080
ENV UPD_STORAGEDIR=/data
ENV UPD_DBFILE=/data/bbolt.db
ENV UPD_DEBUG=1

USER 1001:1001
EXPOSE 8080
VOLUME /data
CMD ["/app/upd"]
