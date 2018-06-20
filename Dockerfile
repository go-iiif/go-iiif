# copied from https://github.com/felixbuenemann/vips-alpine/blob/master/Dockerfile
FROM golang:alpine

ARG VIPS_VERSION=8.6.4

ADD . /go-iiif

RUN wget -O- https://github.com/jcupitt/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz | tar xzC /tmp \
    && apk update \
    && apk upgrade \
    && apk add \
    make libc-dev gcc g++ \
    zlib libxml2 glib gobject-introspection \
    libjpeg-turbo libexif lcms2 fftw giflib libpng \
    libwebp orc tiff poppler-glib librsvg libgsf openexr \
    && apk add --virtual vips-dependencies build-base \
    zlib-dev libxml2-dev glib-dev gobject-introspection-dev \
    libjpeg-turbo-dev libexif-dev lcms2-dev fftw-dev giflib-dev libpng-dev \
    libwebp-dev orc-dev tiff-dev poppler-dev librsvg-dev libgsf-dev openexr-dev \
    && cd /tmp/vips-${VIPS_VERSION} \
    && ./configure --prefix=/usr \
       		   --libdir=/usr/lib64 \
                   --disable-static \
                   --disable-dependency-tracking \
                   --enable-silent-rules \
    && make -s install-strip \
    && cd $OLDPWD \
    && rm -rf /tmp/vips-${VIPS_VERSION} \
    && apk del --purge vips-dependencies \
    && rm -rf /var/cache/apk/* \
    && cd /go-iiif \
    && make bin

# pkg-config --cflags vips vips vips vips
# Package vips was not found in the pkg-config search path.
# Perhaps you should add the directory containing `vips.pc'
# to the PKG_CONFIG_PATH environment variable
# Package 'vips', required by 'virtual:world', not found
# Package 'vips', required by 'virtual:world', not found
# Package 'vips', required by 'virtual:world', not found
# Package 'vips', required by 'virtual:world', not found
# pkg-config: exit status 1
# make: *** [Makefile:66: bin] Error 2

# Step 4/8 : RUN pkg-config --variable pc_path pkg-config
#  ---> Running in 01b282caea91
# /usr/local/lib/pkgconfig:/usr/local/share/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig
# RUN pkg-config --variable pc_path pkg-config

COPY /go-iiif/bin/iiif-server /bin/iiif-server

EXPOSE 8080

# RUN ME...
