FROM golang:1.8-alpine

RUN mkdir -p /go/src/github.com/itglobal/dashboard
WORKDIR /go/src/github.com/itglobal/dashboard
COPY . .
RUN go get
RUN go build
ENTRYPOINT "./dashboard"