FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o nome_teu_servico
 
 
FROM scratch
WORKDIR /app
COPY --from=builder /app/nome_teu_servico /app/nome_teu_servico
 
CMD ["/app/nome_teu_servico"]