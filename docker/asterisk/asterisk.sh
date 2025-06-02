#!/bin/sh

set -ex

# generate certificates
sh ./cert-generate.sh

openrc
rc-service asterisk start &

sleep infinity