FROM golang:1.21 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

CMD ["./main"]