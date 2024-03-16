FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api-hokkaido
 
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-hokkaido .


FROM scratch
WORKDIR /app
COPY --from=builder /app/api-hokkaido /app/api-hokkaido
 
CMD ["/app/api-hokkaido"]