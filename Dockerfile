FROM golang:1-buster AS builder
ENV GO111MODULE on
RUN mkdir /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /src
WORKDIR /src
RUN make setup
RUN make build

FROM alpine:3.11
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=builder /src/cli-template /usr/local/bin
ENTRYPOINT ["/usr/local/bin/cli-template", "serve"]
