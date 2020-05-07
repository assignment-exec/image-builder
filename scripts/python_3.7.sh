#!/bin/bash

# Installation commands for python 3.7
apt-get update
apt install software-properties-common
add-apt-repository ppa:deadsnakes/ppa
apt install python3.7
alias python=python3
apt-get update
apt install python3-pip
