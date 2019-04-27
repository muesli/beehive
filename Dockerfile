FROM golang:alpine AS builder

LABEL authors="Gabriel Alacchi: alacchi.g@gmail.com, Christian Muehlhaeuser: muesli@gmail.com"

# Install git & make
# Git is required for fetching the dependencies
RUN apk update && \
    apk add --no-cache git make ca-certificates && \
    update-ca-certificates

# Set the working directory for the container
WORKDIR /go/beehive

# Build the binary
COPY . .
RUN make embed

FROM alpine

RUN apk update && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

COPY --from=builder /go/beehive/beehive /go/bin/beehive

# Where the admin interface will be served from
ENV CANONICAL_URL=http://localhost:8181

# Expose the application port
EXPOSE 8181

# create a volume for the configuration persistence
VOLUME /conf

# This form of ENTRYPOINT allows the beehive process to catch signals from the `docker stop` command
ENTRYPOINT /go/bin/beehive -config /conf/beehive.conf -bind 0.0.0.0:8181 -canonicalurl ${CANONICAL_URL}
