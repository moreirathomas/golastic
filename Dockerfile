FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/moreirathomas/golastic
COPY . .
RUN go get -d -v ./...
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/golastic cmd/main.go

FROM alpine:latest
WORKDIR /root
COPY .env.docker .
COPY --from=builder /go/src/github.com/moreirathomas/golastic/bin/golastic .
EXPOSE 9999

CMD ["./golastic", "-env-file", "./.env.docker"]  
