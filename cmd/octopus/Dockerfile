FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY .. .
RUN CGO_ENABLED=0 go build -o octopus -mod=vendor -ldflags='-s -w'  -installsuffix cgo cmd/octopus/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /octopus
COPY --from=builder ./build/octopus ./cmd/


EXPOSE 8080

ENTRYPOINT ["./cmd/octopus","-config=/configs/octopus.yml"]
