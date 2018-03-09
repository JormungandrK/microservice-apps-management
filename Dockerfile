### Multi-stage build
FROM jormungandrk/goa-build as build

COPY . /go/src/github.com/JormungandrK/microservice-apps-management
RUN go install github.com/JormungandrK/microservice-apps-management


### Main
FROM alpine:3.7

COPY --from=build /go/bin/microservice-apps-management /usr/local/bin/microservice-apps-management
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-apps-management"]
