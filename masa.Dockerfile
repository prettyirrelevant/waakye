# =============== build stage ===============
FROM golang:1.19.6-buster AS build

WORKDIR /app

ENV BUILD_DIR=/app/bin

COPY go.* ./

RUN go mod download -x all

COPY . .

RUN CGO_ENABLED=0 go build -o masa-app -v cmd/masa/main.go


# =============== final stage ===============
FROM chromedp/headless-shell AS final

COPY --from=build /app/masa-app .

EXPOSE 5001

ENTRYPOINT [ "./masa-app" ]
