#!/bin/sh

apt-get update
apt-get upgrade -y
apt-get install -y git htop sysstat ufw fail2ban unattended-upgrades unzip
dpkg-reconfigure -f noninteractive --priority=low unattended-upgrades

# YMMV - adjust to taste...
apt-get install -y emacs24-nox 

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

# install go

VERSION="1.7.3"
DIST="go${VERSION}.linux-amd64.tar.gz"
SOURCE="https://storage.googleapis.com/golang/${DIST}"

HASH="508028aac0654e993564b6e2014bf2d4a9751e3b286661b0b0040046cf18028e"

if [ ! -d /usr/local/go${VERSION} ]
then
    cd /tmp

    wget ${SOURCE}

    FNAME=`basename ${SOURCE}`
    SRC_HASH=`shasum -a 256 /tmp/${FNAME} | awk '{ print $1 }'`

    if [ "${SRC_HASH}" != "${HASH}" ]
    then
	echo "Weird hash (${SRC_HASH}), expected ${HASH}"
	exit 1
    fi

    tar -xvzf ${DIST}

    if [ -f /tmp/${DIST} ]
    then
	rm /tmp/${DIST}
    fi

    mv /tmp/go /usr/local/go${VERSION}

    if [ -L /usr/local/go ]
    then
	rm /usr/local/go
    fi

    ln -s /usr/local/go${VERSION} /usr/local/go

    for BIN in go gofmt godoc
    do
	
	if [ -L /usr/local/bin/${BIN} ]
	        then
	            rm /usr/local/bin/${BIN}
		    fi

	ln -s /usr/local/go/bin/${BIN} /usr/local/bin/${BIN}
    done

    cd -
fi
