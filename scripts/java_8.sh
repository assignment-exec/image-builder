#!/bin/bash

# Installation commands for java 8
apt update -y
apt install -y apt-transport-https
apt install -y ca-certificates
apt install -y wget dirmngr gnupg
apt install -y software-properties-common
# Fetching the key.
wget -qO - https://adoptopenjdk.jfrog.io/adoptopenjdk/api/gpg/key/public | apt-key add -
# Adding adoptopenjdk repository now that we have the key.
add-apt-repository --yes https://adoptopenjdk.jfrog.io/adoptopenjdk/deb/
# update sources and install java 8 jdk.
apt update -y
apt install -y adoptopenjdk-8-hotspot