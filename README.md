# 🚀 Go AI Stream - Polyglot Microservice Architecture

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Python Version](https://img.shields.io/badge/Python-3.11+-3776AB?style=for-the-badge&logo=python&logoColor=white)
![Rust Version](https://img.shields.io/badge/Rust-2021-000000?style=for-the-badge&logo=rust&logoColor=white)
![CI/CD Status](https://img.shields.io/badge/CI%2FCD-Passing-brightgreen?style=for-the-badge&logo=githubactions&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

A production-ready, high-performance polyglot microservice backend engineered for real-time AI streaming responses, low-latency processing, and robust session management.

---

## 🏗 System Architecture

The project leverages the unique strengths of three programming environments:

flowchart TD
    Client["Client / Frontend"] -->|HTTP / SSE| Gateway["Go API Gateway (:8080)<br/>(Rate Limit, Auth, CI)"]
    Gateway -->|gRPC / HTTP| PythonWorker["Python AI Worker<br/>(LLM Engine / Inference)"]
    Gateway -->|gRPC / HTTP| RustWorker["Rust Performance Worker<br/>(High-Speed Data Engine)"]

Service Roles:
- Go API Gateway (:8080): Entry point for all client communication. Features concurrency-safe request handling, custom RateLimiter, panic Recovery middleware, and chunked HTTP response streaming.
- Python AI Engine: Handles heavy LLM/AI prompt orchestration and dynamic inference workflows.
- Rust Worker: Executes high-throughput data transformations and computationally intensive operations with memory safety.

---

## ✨ Key Features

- Real-time Streaming Response: Native SSE/Chunked HTTP transfer support for interactive AI chat interface.
- Session Continuity: Thread-safe MemoryStore mapping for tracking stateful conversation context across unique session IDs.
- Resilient Infrastructure: Built-in HTTP middleware preventing server crashes (Recovery) and mitigating abuse (RateLimiter).
- Automated Quality Assurance: 100% green CI/CD pipeline integrated via GitHub Actions covering linting and comprehensive package unit tests.

---

## 🛠 Tech Stack

| Domain | Technology |
| :--- | :--- |
| API Gateway | Go 1.22, Standard net/http |
| AI Worker | Python 3.11, AsyncIO |
| Performance Engine | Rust 2021 Edition |
| CI/CD & DevOps | GitHub Actions, Docker, Docker Compose |
| Testing | Go Standard testing Package |

---

## ⚡ Quick Start

### Prerequisites
- Docker & Docker Compose installed.
- Or Go 1.22+ installed locally.

### Option 1: Run via Docker Compose (Recommended)

git clone [https://github.com/YOUR_USERNAME/go-ai-stream.git](https://github.com/YOUR_USERNAME/go-ai-stream.git)
cd go-ai-stream
docker-compose up --build

Access the application web interface at http://localhost:8080.

### Option 2: Run Go Gateway Locally

go run cmd/server/main.go

---

## 🧪 Testing

Execute automated unit tests across all internal Go packages (middleware, storage, service, handlers):

go test -v ./...

---

## 📄 License

This project is open-source and available under the MIT License.
