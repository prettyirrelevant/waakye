FROM golang:1.19.6-buster AS build

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh && \
    curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin/

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o waakye-app -v cmd/waakye/main.go

EXPOSE 8000

CMD ["./waakye-app"]
