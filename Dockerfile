FROM golang:1.7.1-wheezy

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/tleyden/cecil

RUN go get -d -t github.com/tleyden/cecil/... && \
    go install github.com/tleyden/cecil/...

ENTRYPOINT /go/bin/cecil
