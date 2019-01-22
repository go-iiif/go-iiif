# cobbled together from
# https://github.com/felixbuenemann/vips-alpine/blob/master/Dockerfile
# https://github.com/mikestead/docker-imaginary-alpine/blob/master/Dockerfile

FROM golang:alpine as builder

ARG VIPS_VERSION=8.7.4

ENV VIPS_DIR=/vips
ENV PKG_CONFIG_PATH=${VIPS_DIR}/lib/pkgconfig:$PKG_CONFIG_PATH

ADD . /go-iiif

RUN wget -O- https://github.com/libvips/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz | tar xzC /tmp \
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
    libwebp orc tiff poppler-glib librsvg libgsf openexr \
    ca-certificates

RUN mkdir /etc/iiif-server
RUN mkdir /usr/local/iiif-server

VOLUME /etc/iiif-server
VOLUME /usr/local/iiif-server

EXPOSE 8080

COPY cmd/docker-entrypoint.sh /bin/entrypoint.sh

ENTRYPOINT "/bin/entrypoint.sh"

# This is just here as a "for example" and because I can never remember the syntax for this stuff - it
# assumes that you've created the /usr/local/go-iiif/docker/... directories (for mounting with the volumes
# above) locally and copied in the relevant config file and images; adjust to taste (20180620/thisisaaronland)
#
# > docker run -it -p 6161:8080 -e IIIF_SERVER_CONFIG=/etc/iiif-server/config.json -v /usr/local/go-iiif/docker/etc:/etc/iiif-server -v /usr/local/go-iiif/docker/images:/usr/local/iiif-server go-iiif
# 2018/06/20 23:03:10 Listening for requests at 0.0.0.0:8080
#
# curl localhost:6161/test2.jpg/info.json
# {"@context":"http://iiif.io/api/image/2/context.json","@id":"http://localhost:6161/test.jpg","@type":"iiif:Image","protocol":"http://iiif.io/api/image","width":3897,"height":4096,"profile":["http://iiif.io/api/image/2/level2.json",{"formats":["gif","webp","jpg","png","tif"],"qualities":["default","color","dither"],"supports":["full","regionByPx","regionByPct","regionSquare","sizeByDistortedWh","sizeByWh","full","max","sizeByW","sizeByH","sizeByPct","sizeByConfinedWh","none","rotationBy90s","mirroring","noAutoRotate","baseUriRedirect","cors","jsonldMediaType"]}],"service":[{"@context":"x-urn:service:go-iiif#palette","profile":"x-urn:service:go-iiif#palette","label":"x-urn:service:go-iiif#palette","palette":[{"name":"#2f2013","hex":"#2f2013","reference":"vibrant"},{"name":"#9e8e65","hex":"#9e8e65","reference":"vibrant"},{"name":"#c6bca6","hex":"#c6bca6","reference":"vibrant"},{"name":"#5f4d32","hex":"#5f4d32","reference":"vibrant"}]}]}
