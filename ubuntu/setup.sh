#!/bin/sh

apt-get update
apt-get upgrade -y
apt-get install -y git htop sysstat ufw fail2ban unattended-upgrades unzip
dpkg-reconfigure -f noninteractive --priority=low unattended-upgrades

# YMMV - adjust to taste...
apt-get install -y emacs24-nox 

VERSION=`cat /etc/os-release | grep VERSION_ID | awk -F '=' '{ print $2 }'`

if [ "${VERSION}" = "\"14.04\"" ]
then
    apt-get install build-essential pkg-config glib2.0-dev libxml2-dev
    wget http://www.vips.ecs.soton.ac.uk/supported/current/vips-8.3.3.tar.gz
    tar -xvzf vips-8.3.3.tar.gz
    cd vips-8.3.3/
    ./configure
    make
    make install
    ldconfig
    cd -
    rm -rf vips-8.3.3
else
    apt-get install -y libvips-dev
fi

# install go

VERSION="1.7"
DIST="go${VERSION}.linux-amd64.tar.gz"
SOURCE="https://storage.googleapis.com/golang/${DIST}"

HASH="702ad90f705365227e902b42d91dd1a40e48ca7f67a2f4b2fd052aaa4295cd95"

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

# install go-iiif

git clone https://github.com/thisisaaronland/iiif.git
cd iiif
export GOPATH=`pwd`
/usr/local/bin/go get "github.com/gorilla/mux"
/usr/local/bin/go get "gopkg.in/h2non/bimg.v1"
/usr/local/bin/go build
cp ./iiif /usr/local/bin/iiif
cd -
chown -R ubuntu.ubuntu iiif
