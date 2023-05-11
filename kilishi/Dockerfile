# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.20-bullseye as build

WORKDIR /opt/app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 go build -o /opt/app/kilishi ./cmd/kilishi/main.go


# -----------------------------------------------------------------------------
#  Final Stage
# -----------------------------------------------------------------------------
FROM alpine:latest as final

WORKDIR /opt/app

COPY --from=build /opt/app/kilishi /opt/app/kilishi

RUN apk -U upgrade && \
    apk add --no-cache dumb-init ca-certificates

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

EXPOSE 8000

CMD ["./kilishi"]
