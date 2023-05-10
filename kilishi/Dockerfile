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

RUN go build -o /opt/app/kilishi ./cmd/kilishi/main.go


# -----------------------------------------------------------------------------
#  Final Stage
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian11 as final

WORKDIR /opt/app

COPY --from=build /opt/app/kilishi /opt/app/kilishi

EXPOSE 8000

CMD ["./kilishi"]
