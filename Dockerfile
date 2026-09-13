# Giai đoạn 1: Build ứng dụng
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Giai đoạn 2: Chạy ứng dụng nhẹ nhàng trên Alpine
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /server .
EXPOSE 8080
CMD ["./server"]
