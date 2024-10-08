FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -v ./cmd/service && go build -v ./cmd/kafka-consumer