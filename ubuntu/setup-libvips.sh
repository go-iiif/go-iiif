#!/bin/sh

apt-get install -y build-essential pkg-config glib2.0-dev libxml2-dev libjpeg-dev libpng-dev libgif-dev libwebp-dev libtiff-dev libmagick-dev librsvg2-dev

VERSION=`cat /etc/os-release | grep VERSION_ID | awk -F '=' '{ print $2 }'`

# libvips 8.4 has been released and bimg has been updated to use it but there
# appear to still be memory/pointer release errors (20161001/thisisaaronland)

VIPS_MAJOR='8.4'
VIPS_VERSION='8.4.2'

# See this: At the moment it seems easier and more reliable to install from
# source. One day it will all install from apt... (20160930/thisisaaronland)

# if [ "${VERSION}" = "\"14.04\"" ]
# then

    wget http://www.vips.ecs.soton.ac.uk/supported/${VIPS_MAJOR}/vips-${VIPS_VERSION}.tar.gz
    tar -xvzf vips-${VIPS_VERSION}.tar.gz
    cd vips-${VIPS_VERSION}/
    ./configure
    make
    make install
    ldconfig
    cd -
    rm -rf vips-${VIPS_VERSION}
    rm -rf vips-${VIPS_VERSION}.tar.gz

# else
#     apt-get install -y libjpeg-dev libpng-dev libgif-dev libwebp-dev libtiff-dev libmagickcore-dev librsvg2-dev libvips-dev
# fi

