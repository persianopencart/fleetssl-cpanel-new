FROM golang:1.20-alpine
MAINTAINER FleetSSL cPanel <support@fleetssl.com>

RUN apk update && \
        apk add curl \
        bash curl-dev ruby-dev build-base ruby ruby-io-console ruby-bundler \
        libffi libffi-dev gawk rpm gcc libc-dev cpio tar

RUN gem install fpm

RUN mkdir -p /build

WORKDIR /build

CMD make clean rpm
