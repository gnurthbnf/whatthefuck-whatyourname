import asyncio
import logging
import os
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

# Cấu hình logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("ai-worker")

app = FastAPI(
    title="AI Worker Service",
    description="Python microservice xử lý AI streaming tích hợp cùng Go Gateway",
    version="1.0.0"
)

# Cấu hình CORS cho phép gọi xuyên dịch vụ
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Khung dữ liệu đầu vào chuẩn Pydantic
class PromptRequest(BaseModel):
    prompt: str
    max_tokens: int = 100
    temperature: float = 0.7

@app.get("/health")
async def health_check():
    """Kiểm tra trạng thái sống của service"""
    return {"status": "healthy", "service": "python-ai-worker"}

async def generate_ai_stream(prompt: str):
    """
    Hàm giả lập quá trình AI sinh tokens theo dạng Streaming (SSE).
    Sau này cậu có thể thay thế phần này bằng gọi Ollama/OpenAI SDK trực tiếp.
    """
    logger.info(f"Đang xử lý prompt từ Go server: {prompt}")
    
    # Chuỗi phản hồi mẫu chia nhỏ thành từng token
    response_tokens = [
        "Xin ", "chào ", "cậu! ", "Hệ ", "thống ", "AI ", 
        "Python ", "đã ", "nhận ", "được ", "yêu cầu ", "và ", 
        "đang ", "xử lý ", "dữ liệu ", "thời gian thực ", "cho cậu đây. 🚀"
    ]
    
    for token in response_tokens:
        # Định dạng chuẩn Server-Sent Events (SSE)
        yield f"data: {token}\n\n"
        await asyncio.sleep(0.05)  # Giả lập độ trễ sinh token của AI

@app.post("/api/v1/generate")
async def handle_ai_generation(request: PromptRequest):
    """API chính nhận yêu cầu từ Go Gateway và trả về stream"""
    if not request.prompt.strip():
        raise HTTPException(status_code=400, detail="Prompt không được để trống!")
    
    return StreamingResponse(
        generate_ai_stream(request.prompt), 
        media_type="text/event-stream"
    )
