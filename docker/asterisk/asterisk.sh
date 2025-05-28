#!/bin/sh

set -ex

openrc
rc-service asterisk start &

sleep infinity