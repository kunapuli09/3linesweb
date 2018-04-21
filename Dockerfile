FROM golang:latest

WORKDIR /go/src/github.com/kunapuli09/3linesweb

RUN go get -u github.com/golang/dep/cmd/dep
# Add project directory to Docker image.
ADD . /go/src/github.com/kunapuli09/3linesweb

ENV USER krishnakunapuli
ENV verbose "true"
ENV HTTP_ADDR "0.0.0.0:8888"
ENV HTTP_DRAIN_INTERVAL "1s"
ENV COOKIE_SECRET "z5mOYQcyv3KQHe3W"
ENV DB_URL "root:kk@starpath@tcp(172.17.0.2:3306)/3linesweb?parseTime=true"

#Security settings
ENV HTTP_CERT_FILE ""
ENV HTTP_KEY_FILE ""

#Mail settings
ENV EMAIL_SENDER_ID "krishna.3lines@gmail.com"
ENV EMAIL_RECEIVER_ID "krishna.kunapuli@threelines.us"
ENV AUTH_PASSWORD "meyt2hn32hmy9sx6es"
ENV SMTP_HOST "smtp.gmail.com"
ENV SMTP_PORT "465"


# Replace this with actual PostgreSQL DSN.
ENV DSN $GO_BOOTSTRAP_MYSQL_DSN

RUN dep ensure && dep status

RUN go build

EXPOSE 8888
CMD ["./3linesweb"]
