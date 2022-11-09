FROM golang:1.18-alpine3.16 as build

RUN apk add --update --no-cache ca-certificates make

WORKDIR /go/upctl/
COPY . .
RUN make build-dockerised


FROM alpine:3.16

LABEL org.label-schema.vcs-url="https://github.com/UpCloudLtd/upcloud-cli"

RUN apk add --update --no-cache ca-certificates

COPY --from=build /go/upctl/bin/*dockerised-linux-amd64 /bin/upctl

ENTRYPOINT ["/bin/upctl"]
