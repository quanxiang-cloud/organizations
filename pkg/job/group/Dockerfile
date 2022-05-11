FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o pushdata -mod=vendor -ldflags='-s -w'  -installsuffix cgo pkg/job/group/pushdata.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /pushdata
COPY --from=builder ./build/pushdata ./cmd/

EXPOSE 80

ENTRYPOINT ["./cmd/pushdata","-config=/configs/config.yml"]
