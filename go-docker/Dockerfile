FROM golang:latest

LABEL maintainer="Ben Gavan <email@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
