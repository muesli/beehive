FROM golang:1.12-alpine

LABEL authors="Gabriel Alacchi: alacchi.g@gmail.com, Christian Muehlhaeuser: muesli@gmail.com"

ENV CANONICAL_URL=http://localhost:8181

# Expose the application port
EXPOSE 8181

# Install beehive
RUN apk add -U --no-cache git make && \
    git clone https://github.com/muesli/beehive.git && \
    cd beehive && \
    make

# Set the working directory for the container
WORKDIR /go/beehive

# create a volume for the configuration persistence
VOLUME /conf

# This form of ENTRYPOINT allows the beehive process to catch signals from the `docker stop` command
ENTRYPOINT ./beehive -config /conf/beehive.conf -bind 0.0.0.0:8181 -canonicalurl ${CANONICAL_URL}
