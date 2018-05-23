#!/bin/sh

apt-get update
apt-get upgrade -y
apt-get install -y git htop sysstat ufw fail2ban unattended-upgrades unzip
dpkg-reconfigure -f noninteractive --priority=low unattended-upgrades

# YMMV - adjust to taste...
apt-get install -y emacs24-nox 
