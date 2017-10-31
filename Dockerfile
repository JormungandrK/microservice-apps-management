### Multi-stage build
FROM golang:1.8.3-alpine3.6 as build

RUN apk --no-cache add git curl openssh

RUN go get -u -v github.com/goadesign/goa/... && \
    go get -u -v github.com/asaskevich/govalidator && \
    go get -u -v gopkg.in/mgo.v2 && \
    go get -u -v github.com/JormungandrK/microservice-tools && \
    go get -u -v github.com/JormungandrK/microservice-security/...

COPY . /go/src/github.com/JormungandrK/microservice-apps-management
RUN go install github.com/JormungandrK/microservice-apps-management


### Main
FROM alpine:3.6

COPY --from=build /go/bin/microservice-apps-management /usr/local/bin/microservice-apps-management
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-apps-management"]
