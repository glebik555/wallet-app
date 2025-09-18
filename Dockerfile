FROM golang:1.23-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git

COPY . .

RUN go mod tidy && go mod download

RUN go test ./... -v

RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-app ./cmd/app

FROM scratch

COPY --from=builder /wallet-app /wallet-app
COPY .env /.env

EXPOSE 8080

ENTRYPOINT ["/wallet-app"]
