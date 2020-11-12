FROM golang:1.15-alpine as builder

RUN apk add --no-cache git alpine-sdk

ADD . /code
WORKDIR /code

# build the source
RUN go build -o spider cmd/main.go

# use a minimal alpine image
FROM alpine:3.12

# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# set working directory
WORKDIR /go/bin

COPY --from=builder /code/spider .

USER 1001
# run the binary
CMD ["./spider"]