# Giai đoạn build
FROM rust:1.77-alpine AS builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Giai đoạn chạy thực tế
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/target/release/rust_worker .
EXPOSE 8001
CMD ["./rust_worker"]
