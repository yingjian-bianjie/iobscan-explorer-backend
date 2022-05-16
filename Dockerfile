FROM golang:1.16-alpine3.14 as builder

# Set up dependencies
ENV PACKAGES make git libc-dev bash
WORKDIR $GOPATH/src
COPY . .

# Install minimum necessary dependencies, build binary
RUN apk add --no-cache $PACKAGES && make install

FROM alpine:3.12
COPY --from=builder /go/bin/ddcparser /usr/local/bin/ddcparser
RUN mkdir $HOME/.ddc-parser
CMD ddcparser start