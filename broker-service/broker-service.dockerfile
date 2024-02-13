FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o brockerApp ./cmd/api

RUN chmod +x /app/brockerApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brockerApp /app

CMD [ "/app/brockerApp" ]