# Go AI Stream

Hệ thống Microservice xử lý AI Streaming real-time hiệu suất cao viết bằng ngôn ngữ Go, hỗ trợ lưu trữ phiên làm việc qua Redis/RAM và cơ chế Fallback thông minh giữa các AI Provider.

## Cấu trúc thư mục
- `cmd/server/`: Điểm khởi chạy ứng dụng (Entry point).
- `internal/config/`: Quản lý biến môi trường và thiết lập hệ thống.
- `internal/handlers/`: Xử lý HTTP request, SSE streaming và giao diện web.
- `internal/middleware/`: Bộ lọc chống spam (Rate limiter).
- `internal/service/`: Logic cốt lõi tích hợp cơ chế gọi AI dự phòng (Fallback).
- `internal/storage/`: Tích hợp lưu trữ phiên qua In-Memory RAM hoặc Redis phân tán.

## Hướng dẫn chạy nhanh bằng Docker
```bash
docker-compose up --build
