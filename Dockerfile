FROM golang:1.21

RUN mkdir /app
WORKDIR /app

COPY . .
COPY .env .

RUN go mod download

RUN go build -o /build cmd/server/main.go

CMD ["/build"]