# build the container - this is used to cache go.mod files, otherwise it is too slow to re-download every time
FROM golang:alpine as base

WORKDIR /opt/

COPY go.mod /opt/
COPY go.sum /opt/

RUN GOOS=linux GOARCH=amd64 go mod download

# build the container - this builds my code
FROM golang:alpine as builder
COPY --from=base /go/pkg/mod/ /go/pkg/mod/
COPY --from=base /opt/go.mod /opt/go.mod
COPY --from=base /opt/go.sum /opt/go.sum

WORKDIR /opt/
COPY src/ /opt/src/

RUN cd /opt/src && \
    GOOS=linux GOARCH=amd64 go build -o /opt/service 

# runtime
FROM alpine

COPY --from=builder /opt/service /opt/service

EXPOSE 3000

ENV SERVICE="dio-scope3"
ENV ENV="LOCAL"
ENV VERSION="1"
ENV PORT="3000"
ENV SCOPE3_API_TOKEN=""

CMD ["/opt/service"]