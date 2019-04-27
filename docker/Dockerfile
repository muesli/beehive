FROM golang:1.11-alpine

RUN apk add -U --no-cache curl git make g++ && \
    curl "https://raw.githubusercontent.com/golang/dep/master/install.sh" | sh && \
    go get -d -v github.com/muesli/beehive && \
    cd ${GOPATH}/src/github.com/muesli/beehive && \
    make get-deps && \
    make


FROM alpine:3.8

RUN apk add -U --no-cache ca-certificates

WORKDIR /app

COPY --from=0 /go/src/github.com/muesli/beehive/beehive .

ENTRYPOINT ["./beehive"]
CMD ["-bind", "0.0.0.0:8181", "-canonicalurl", "http://localhost:8181"]
