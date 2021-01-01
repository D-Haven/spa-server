FROM golang:1.15 AS builder
ARG version=1.0.0

WORKDIR /go/src/github.com/D-Haven/spa-server/
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV RELEASE=${version}

RUN go get -d -v
RUN PROJECT=d-haven.org/spa-server && \
    COMMIT=$(git rev-parse --short HEAD) && \
    BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%S') && \
    go build -ldflags "-X ${PROJECT}/version.Release=${RELEASE} \
                -X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
             	-a -o /go/bin/spa-server

RUN go test ./... -cover

##########################################

FROM scratch

COPY --from=builder /go/bin/spa-server /

EXPOSE 8443/tcp
ENTRYPOINT ["/spa-server"]