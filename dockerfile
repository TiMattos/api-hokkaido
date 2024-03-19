FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-hokkaido

FROM scratch
WORKDIR /app
COPY --from=builder /app/api-hokkaido /app/api-hokkaido
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/app/api-hokkaido"]
