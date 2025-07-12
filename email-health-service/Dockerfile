# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN cd cmd && go build -o /email-health-service

FROM alpine:latest
COPY --from=builder /email-health-service /email-health-service
EXPOSE 8081
CMD ["/email-health-service"]