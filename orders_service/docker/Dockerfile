FROM golang:1.18-alpine as builder

WORKDIR /root

COPY ./orders_service/. ./
RUN go mod download && \
    go mod verify && \
    go build -o ./app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /root/app ./app

EXPOSE 8085

CMD ["./app"]