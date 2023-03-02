FROM golang:1.18-alpine3.16 as build

RUN apk add --update --no-cache ca-certificates git make

WORKDIR /go/upctl/
COPY . .
RUN make build-dockerised


FROM alpine:3.16

RUN apk add --update --no-cache ca-certificates jq
COPY --from=build /go/upctl/bin/*dockerised-linux-amd64 /bin/upctl

ENTRYPOINT ["/bin/upctl"]
