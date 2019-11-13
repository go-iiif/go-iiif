FROM golang:1.12-alpine as builder

ADD . /go-iiif

RUN apk update \
    && apk upgrade \
    && apk add make \
    && cd /go-iiif \
    && make cli-tools

FROM alpine

COPY --from=builder /go-iiif/bin/iiif-process /bin/iiif-process
COPY --from=builder /go-iiif/bin/iiif-server /bin/iiif-server
COPY --from=builder /go-iiif/bin/iiif-tile-seed /bin/iiif-tile-seed

RUN apk update \
    && apk upgrade \
    && apk add \    
    ca-certificates
    
RUN mkdir /etc/go-iiif
RUN mkdir /usr/local/go-iiif

VOLUME /etc/go-iiif
VOLUME /usr/local/go-iiif