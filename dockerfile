FROM golang:latest

WORKDIR /app

RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/data

RUN go test -v ./tasks

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]