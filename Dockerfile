# build the container
FROM golang:alpine as builder

WORKDIR /opt/

COPY . /opt/

RUN GOOS=linux GOARCH=amd64 go build -o /opt/service /opt/src/


# runtime
FROM alpine

COPY --from=builder /opt/service /opt/service

EXPOSE 3000

ENV SERVICE="dio-scope3"
ENV ENVIRONMENT="LOCAL"
ENV VERSION="1"