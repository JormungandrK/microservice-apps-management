### Multi-stage build
FROM golang:1.17.3-alpine3.15 as build

RUN apk --no-cache add git curl openssh

COPY . /go/src/github.com/Microkubes/microservice-apps-management

RUN cd /go/src/github.com/Microkubes/microservice-apps-management && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install

### Main
FROM alpine:3.15

COPY --from=build /go/src/github.com/Microkubes/microservice-apps-management/config.json /config.json
COPY --from=build /go/bin/microservice-apps-management /microservice-apps-management

EXPOSE 8080

CMD ["/microservice-apps-management"]
