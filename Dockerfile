FROM golang:alpine

WORKDIR /go/src/blockcredit

RUN apk update && apk add git

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000-4000

WORKDIR /go/bin

ENTRYPOINT [ "./blockcredit"]

CMD [ "-p 3000" ]