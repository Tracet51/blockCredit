FROM golang:alpine

WORKDIR /go/src/blockcredit

RUN apk update && apk add git

COPY *.go ./

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000-9000

WORKDIR /go/bin

ENTRYPOINT [ "./blockcredit", "-p"]