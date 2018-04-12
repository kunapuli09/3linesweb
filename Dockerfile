FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/kunapuli09/3linesweb

ENV USER krishnakunapuli
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET z5mOYQcyv3KQHe3W

# Replace this with actual PostgreSQL DSN.
ENV DSN $GO_BOOTSTRAP_MYSQL_DSN

WORKDIR /go/src/github.com/kunapuli09/3linesweb

RUN godep go build

EXPOSE 8888
CMD ./3linesweb
