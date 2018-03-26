#!/bin/sh
SRC=$(cd "$(dirname "$0")/.."; pwd)
if [ ! -d "${SRC}/tmp" ]; then
  mkdir "${SRC}/tmp"
fi

openssl genrsa -out "${SRC}/tmp/private_key.rsa"
openssl rsa -in "${SRC}/tmp/private_key.rsa" -pubout > "${SRC}/tmp/public_key.rsa.pub"