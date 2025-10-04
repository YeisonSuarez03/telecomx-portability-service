FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rest-service ./cmd/

FROM gcr.io/distroless/base-debian12
WORKDIR /app

ENV KAFKA_BROKERS kafka.railway.internal:29092
ENV KAFKA_CLIENT_ID telecomx-portability-service
ENV KAFKA_GROUP_ID telecomx-portability-consumer
ENV KAFKA_TOPIC Customer
ENV PORT 8080

COPY --from=builder /app/rest-service .

EXPOSE 8080

CMD ["./rest-service"]