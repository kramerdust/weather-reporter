FROM golang:stretch

WORKDIR /app
COPY . .

RUN go install -v ./cmd/service.go

ENTRYPOINT /go/bin/service
EXPOSE 8080