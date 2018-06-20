# cobbled together from
# https://github.com/felixbuenemann/vips-alpine/blob/master/Dockerfile
# https://github.com/mikestead/docker-imaginary-alpine/blob/master/Dockerfile

FROM golang:alpine as builder

ARG VIPS_VERSION=8.6.4

ADD . /go-iiif

ENV VIPS_DIR=/vips
ENV PKG_CONFIG_PATH=${VIPS_DIR}/lib/pkgconfig:$PKG_CONFIG_PATH

RUN wget -O- https://github.com/jcupitt/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz | tar xzC /tmp \
    && apk update \
    && apk upgrade \

    && apk add \
    make libc-dev gcc \
    zlib libxml2 glib gobject-introspection \
    libjpeg-turbo libexif lcms2 fftw giflib libpng \
    libwebp orc tiff poppler-glib librsvg libgsf openexr \

    && apk add --virtual vips-dependencies build-base \
    zlib-dev libxml2-dev glib-dev gobject-introspection-dev \
    libjpeg-turbo-dev libexif-dev lcms2-dev fftw-dev giflib-dev libpng-dev \
    libwebp-dev orc-dev tiff-dev poppler-dev librsvg-dev libgsf-dev openexr-dev \

    && cd /tmp/vips-${VIPS_VERSION} \
    && ./configure --prefix=${VIPS_DIR} \
                   --disable-static \
		   --without-python \
                   --disable-dependency-tracking \
                   --enable-silent-rules \
    && make -s install-strip \

    && cd /go-iiif \
    && make bin

FROM alpine

COPY --from=builder /vips/lib/ /usr/local/lib
COPY --from=builder /go-iiif/bin/iiif-server /bin/iiif-server

RUN apk update \
    && apk upgrade \
    && apk add \
    zlib libxml2 glib gobject-introspection \
    libjpeg-turbo libexif lcms2 fftw giflib libpng \
    libwebp orc tiff poppler-glib librsvg libgsf openexr

RUN mkdir /etc/iiif-server
COPY config.json /etc/iiif-server/config.json

RUN mkdir /usr/local/iiif-server
COPY example/images/184512_5f7f47e5b3c66207_x.jpg /usr/local/iiif-server/184512_5f7f47e5b3c66207_x.jpg

EXPOSE 8080

# RUN ME...

ENTRYPOINT [ "/bin/iiif-server", "-config",  "/etc/iiif-server/config.json" ]
