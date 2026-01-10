FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /golang_twitter

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /golang_twitter/main .

CMD ["./main"]

EXPOSE 8080