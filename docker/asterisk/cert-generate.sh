#! /bin/sh

if [ -z $1  ]
then
    base="certificate"
else
    base=$1
fi

path="/home/asterisk/certs/"
key="$path$base.key"
cert="$path$base.crt"

if [ -f $key ]
then
    echo "$key and $cert exists"
else
    echo "$key and $cert does not exists, creating..."
    password=password
    openssl genrsa -passout pass:$password -out certificate.key.orig

    country="XX"
    state="XX"
    locality="XX"
    organization="XX"
    organizationalunit="XX"
    commonname="XX"
    email="XX"

    openssl req -new -key certificate.key.orig -out certificate.csr -passin pass:$password \
    -subj "/C=$country/ST=$state/L=$locality/O=$organization/OU=$organizationalunit/CN=$commonname/emailAddress=$email"

    openssl rsa -in certificate.key.orig -out certificate.key -passin pass:$password
    openssl x509 -req -days 3650 -in certificate.csr -signkey certificate.key -out certificate.crt

    cp certificate.key $key
    cp certificate.crt $cert
    chown -R asterisk:asterisk $path

    rm certificate.*
fi
