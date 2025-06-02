FROM alpine:3.21 AS local-asterisk
RUN apk update && apk add --no-cache asterisk tini openrc openssl && \
  mkdir -p /run/openrc && touch /run/openrc/softlevel

COPY ./docker/asterisk/certificate/ home/asterisk/certs/
COPY ./docker/asterisk/cert-generate.sh .
RUN chmod +x ./cert-generate.sh
COPY ./docker/asterisk/asterisk.sh .
RUN chmod +x ./asterisk.sh

ENTRYPOINT ["tini", "--", "./asterisk.sh"] 


