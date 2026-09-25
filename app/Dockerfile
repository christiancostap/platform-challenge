FROM golang:1.27.1-alpine3.24 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app .

FROM alpine:3.24

COPY --from=builder /app .

ENTRYPOINT ["/app"]