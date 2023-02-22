FROM golang:1.19.4-bullseye AS builder

ADD . /src
WORKDIR /src
RUN make build

FROM alpine:3.13.5
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /src/bin/main usr/local/bin/main
RUN chmod +x /usr/local/bin/main

ENTRYPOINT ["/usr/local/bin/main"]
