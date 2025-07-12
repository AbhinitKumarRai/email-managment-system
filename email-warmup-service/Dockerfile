# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd && go build -o /email-warmup-service

FROM alpine:latest
COPY --from=builder /email-warmup-service /email-warmup-service
EXPOSE 8080
CMD ["/email-warmup-service"]