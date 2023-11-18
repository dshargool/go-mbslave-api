FROM alpine:3.18 as builder

RUN apk add --no-cache sqlite git make musl-dev build-base

COPY --from=golang:1.21-alpine /usr/local/go/ /usr/local/go/

ENV GOROOT /usr/local/go 
ENV GOPATH /go 
ENV PATH /usr/local/go/bin:$PATH 
ENV CGO_ENABLED 1

RUN mkdir -p /go/app/ /go/bin/
COPY go.mod $GOPATH/app/
COPY go.sum $GOPATH/app/
WORKDIR $GOPATH/app/
RUN go get github.com/mattn/go-sqlite3@v1.14.17

COPY . $GOPATH/app/
RUN go build -o $GOPATH/bin/go-mbslave-api
RUN ldd $GOPATH/bin/go-mbslave-api


FROM scratch
COPY --from=builder /go/bin/go-mbslave-api go-mbslave-api
COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
CMD ["go-mbslave-api"]
