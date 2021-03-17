# https://hub.docker.com/_/golang?tab=tags
FROM golang:1.16.2

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["healthcheck"]