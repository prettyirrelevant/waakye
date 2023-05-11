# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.20-bullseye as build

ENV CGO_ENABLED=0

WORKDIR /opt/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /opt/app/kilishi ./cmd/kilishi/main.go && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 && \
    chmod +x /usr/local/bin/dumb-init


# -----------------------------------------------------------------------------
#  Final Stage
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian11 as final

WORKDIR /opt/app

COPY --from=build /opt/app/kilishi /opt/app/kilishi
COPY --from=build /usr/local/bin/dumb-init /usr/local/bin/dumb-init

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

EXPOSE 8000

CMD ["./kilishi"]
