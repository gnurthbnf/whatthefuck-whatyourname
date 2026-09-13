use axum::{
    routing::get,
    Json, Router,
};
use serde::Serialize;
use std::net::SocketAddr;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    service: String,
}

async fn health_check() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "healthy".to_string(),
        service: "rust-performance-worker".to_string(),
    })
}

#[tokio::main]
async fn main() {
    // Định tuyến API cho service Rust
    let app = Router::new().route("/health", get(health_check));

    // Lắng nghe ở cổng 8001
    let addr = SocketAddr::from(([0, 0, 0, 0], 8001));
    println!("Rust microservice đang chạy tại http://{}", addr);
    
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
