#!/bin/sh

yum upgrade

yum install -y git htop sysstat ufw fail2ban unzip

# YMMV - adjust to taste...
yum install -y emacs-nox 
