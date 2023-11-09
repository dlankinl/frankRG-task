FROM golang:1.21

RUN mkdir /app
WORKDIR /app

COPY . .
COPY .env .

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

RUN go build -o /build cmd/server/main.go

CMD ["/build"]