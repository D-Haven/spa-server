FROM golang:1.15 AS builder

WORKDIR /go/src/github.com/D-Haven/spa-server/
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go get -d -v
RUN go build -a -o /go/bin/spa-server
RUN go test


FROM scratch

COPY --from=builder /go/bin/spa-server /

EXPOSE 8443/tcp
CMD ["/spa-server"]