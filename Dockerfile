
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod download

RUN go build -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssl && \
    mkdir -p /etc/ssl/app && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout /etc/ssl/app/self.key \
      -out    /etc/ssl/app/self.crt \
      -subj "/C=RU/ST=Moscow/L=Moscow/O=Dev/CN=localhost"

WORKDIR /root/

COPY --from=backend-builder /app/main .

CMD ["./main"]