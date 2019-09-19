FROM golang:stretch
RUN apt update

WORKDIR /app
COPY . .

RUN go install -v ./cmd/service.go

ENTRYPOINT /go/bin/service
EXPOSE 3000