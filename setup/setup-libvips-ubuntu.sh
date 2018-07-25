#!/bin/sh
# https://jcupitt.github.io/libvips/install.html

apt-get install -y build-essential pkg-config glib2.0-dev libxml2-dev libjpeg-dev libpng-dev libgif-dev libwebp-dev libtiff-dev libmagick-dev librsvg2-dev

VIPS_VERSION='8.6.5'

# See this: At the moment it seems easier and more reliable to install from
# source. One day it will all install from apt... (20160930/thisisaaronland)

wget https://github.com/jcupitt/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz	
tar -xvzf vips-${VIPS_VERSION}.tar.gz

cd vips-${VIPS_VERSION}/
./configure
make
make install

ldconfig
cd -

rm -rf vips-${VIPS_VERSION}
rm -rf vips-${VIPS_VERSION}.tar.gz
