FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o organizations -mod=vendor -ldflags='-s -w'  -installsuffix cgo cmd/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /organizations
COPY --from=builder ./build/organizations ./cmd/


EXPOSE 80

ENTRYPOINT ["./cmd/organizations","-config=/configs/config.yml"]
