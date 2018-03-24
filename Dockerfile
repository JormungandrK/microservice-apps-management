### Multi-stage build
FROM golang:1.10-alpine3.7 as build

RUN apk --no-cache add git curl openssh

RUN go get -u -v github.com/goadesign/goa/... && \
    go get -u -v github.com/asaskevich/govalidator && \
    go get -u -v github.com/Microkubes/microservice-security/... && \
    go get -u -v github.com/Microkubes/microservice-tools/...


COPY . /go/src/github.com/Microkubes/microservice-apps-management
RUN go install github.com/Microkubes/microservice-apps-management


### Main
FROM alpine:3.7

COPY --from=build /go/bin/microservice-apps-management /microservice-apps-management
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/microservice-apps-management"]
