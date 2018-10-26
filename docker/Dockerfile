FROM golang:1.8

LABEL authors="Gabriel Alacchi: alacchi.g@gmail.com, Christian Muehlhaeuser: muesli@gmail.com"

ENV CANONICAL_URL=http://localhost:8181

# Expose the application port
EXPOSE 8181

# Install beehive
RUN go get github.com/muesli/beehive

# Set the working directory for the container
WORKDIR /go/src/github.com/muesli/beehive

# create a volume for the configuration persistence
VOLUME /conf

# This form of ENTRYPOINT allows the beehive process to catch signals from the `docker stop` command
ENTRYPOINT beehive -config /conf/beehive.conf -bind 0.0.0.0:8181 -canonicalurl ${CANONICAL_URL}
