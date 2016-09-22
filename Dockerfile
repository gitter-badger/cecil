FROM golang:1.7.1-wheezy

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/tleyden/zerocloud

RUN go get -d -t github.com/tleyden/zerocloud/... && \
    go install github.com/tleyden/zerocloud/...

ENTRYPOINT /go/bin/temporary
