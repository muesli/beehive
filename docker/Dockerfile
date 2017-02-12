FROM golang:1.7

# Expose the application port
EXPOSE 8181

# Install beehive
RUN go get github.com/muesli/beehive

# Set the working directory for the container
WORKDIR /go/src/github.com/muesli/beehive

# create a volume for the configuration persistence
VOLUME /conf

# This form of ENTRYPOINT allows the beehive process to catch signals from the `docker stop` command
ENTRYPOINT ["beehive", "-config", "/conf/beehive.conf"]
