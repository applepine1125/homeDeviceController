#!/bin/sh

mkdir ./pkg
cd ./pkg

git clone https://github.com/mjg59/python-broadlink.git
./python-broadlink/setup.py
