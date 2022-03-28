FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o cleandatajobv1 -mod=vendor -ldflags='-s -w'  -installsuffix cgo pkg/job/refactor/cleandatav1.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /cleandatajobv1
COPY --from=builder ./build/cleandatajobv1 ./cmd/

EXPOSE 80

ENTRYPOINT ["./cmd/cleandatajobv1","-config=/configs/config.yml"]
