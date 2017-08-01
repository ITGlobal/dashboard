# Build backend
FROM golang:1.8 as backend
RUN mkdir -p /go/src/github.com/itglobal/dashboard
WORKDIR /go/src/github.com/itglobal/dashboard
COPY . .
RUN go get
RUN go build

# Build frontend
FROM node:8.1 as frontend
RUN mkdir -p /app
WORKDIR /app
COPY /ui /app
RUN yarn install
RUN npm run build:prod

# Application container
FROM alpine:latest
RUN mkdir -p /opt/dashboard
WORKDIR /opt/dashboard
COPY --from=backend /go/src/github.com/itglobal/dashboard/dashboard /opt/dashboard
COPY --from=frontend /app/dist /opt/dashboard/www

ENTRYPOINT "/opt/dashboard"