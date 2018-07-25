#!/bin/sh
# https://jcupitt.github.io/libvips/install.html

yum install -y pkgconfig glib2-devel libpng-devel libjpeg-devel giflib-devel librsvg2-devel poppler-devel libexif-devel expat-devel

VIPS_VERSION='8.6.5'

# See this: At the moment it seems easier and more reliable to install from
# source. One day it will all install from apt... (20160930/thisisaaronland)

wget https://github.com/jcupitt/libvips/releases/download/v${VIPS_VERSION}/vips-${VIPS_VERSION}.tar.gz	
tar -xvzf vips-${VIPS_VERSION}.tar.gz

cd vips-${VIPS_VERSION}/
./configure --lib-dir=/usr/lib64
make
make install

ldconfig
cd -

rm -rf vips-${VIPS_VERSION}
rm -rf vips-${VIPS_VERSION}.tar.gz
