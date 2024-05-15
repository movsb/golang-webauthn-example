#!/bin/bash
#
set -eu

DOMAIN=localhost

openssl genpkey -algorithm rsa > ${DOMAIN}.key
openssl req -new -out ${DOMAIN}.csr -key ${DOMAIN}.key
openssl x509 -req -days 3650 -in ${DOMAIN}.csr -signkey ${DOMAIN}.key -out ${DOMAIN}.crt
