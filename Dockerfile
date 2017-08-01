# Build backend
FROM golang:1.8 as backend
RUN go get -v github.com/gorilla/mux && \
    go get -v github.com/gorilla/websocket && \
    go get -v github.com/kpango/glg && \
    go get -v github.com/mkideal/cli && \
    go get -v github.com/satori/go.uuid && \
    go get -v gopkg.in/mgo.v2 && \
    go get -v github.com/kapitanov/go-teamcity
RUN mkdir -p /go/src/github.com/itglobal/dashboard
WORKDIR /go/src/github.com/itglobal/dashboard
COPY . .
RUN go get -v
RUN go build -v -x

# Build frontend
FROM node:8.1 as frontend
RUN mkdir -p /app
WORKDIR /app
COPY /ui /app
RUN yarn install
RUN npm run build:prod

# Application container
FROM ubuntu:xenial
RUN mkdir -p /app
WORKDIR /app
COPY --from=backend /go/src/github.com/itglobal/dashboard/dashboard /app
COPY --from=frontend /app/dist /app/www
RUN chmod +x /app/dashboard

CMD "/app/dashboard"