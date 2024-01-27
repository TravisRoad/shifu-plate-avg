FROM golang:1.21.6-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o main cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
CMD ["/app/main"]
