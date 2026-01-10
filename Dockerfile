FROM golang:1.25.5-alpine3.23

WORKDIR /golang_twitter

# go installで必要
RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# airを実行（ホットリロード有効）
CMD ["air", "-c", ".air.toml"]

EXPOSE 8080
