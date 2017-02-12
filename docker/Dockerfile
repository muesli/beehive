FROM golang:1.7

# Expose the application port
EXPOSE 8181

# Install beehive
RUN go get github.com/muesli/beehive

# Set the working directory for the container
WORKDIR /go/src/github.com/muesli/beehive

# create a volume for the configuration persistence
VOLUME /conf

# Copy the default config file (users can use their own with docker's -v flag)
COPY beehive.conf /conf/beehive.conf

# This form of ENTRYPOINT allows  the beehive process can catch signals from docker stop or docker kill.
ENTRYPOINT ["beehive", "-config", "/conf/beehive.conf"]
