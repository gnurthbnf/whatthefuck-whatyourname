package handlers

import (
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <title>Go AI Stream - Microservice Architecture</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; max-width: 750px; margin: 40px auto; padding: 25px; background: #f0f2f5; color: #333; }
        .card { background: white; padding: 25px; border-radius: 10px; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
        textarea { width: 100%; height: 100px; padding: 12px; font-size: 15px; border-radius: 6px; border: 1px solid #ccd0d5; box-sizing: border-box; resize: vertical; }
        button { padding: 12px 24px; font-size: 16px; background: #0066cc; color: white; border: none; border-radius: 6px; cursor: pointer; margin-top: 12px; font-weight: bold; transition: background 0.2s; }
        button:hover { background: #004fb3; }
        #output { white-space: pre-wrap; background: #fafbfc; padding: 18px; border: 1px solid #e1e4e8; border-radius: 6px; margin-top: 20px; min-height: 150px; line-height: 1.6; }
        .meta { font-size: 13px; color: #65676b; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="card">
        <h2>Go AI Stream Production Service</h2>
        <textarea id="prompt" placeholder="Nhập câu hỏi hoặc yêu cầu cho AI...">Cậu hãy giải thích kiến trúc Clean Architecture trong Go bằng vài ý chính.</textarea>
        <br>
        <button onclick="startStream()">Gửi Streaming Request</button>
        <div class="meta" id="status">Trạng thái: Sẵn sàng | Session ID: <span id="sid"></span></div>
        <h3>Phản hồi từ AI (Real-time Stream):</h3>
        <div id="output"></div>
    </div>

    <script>
        let sessionId = localStorage.getItem('chat_session_id');
        if (!sessionId) {
            sessionId = 'session_' + Math.random().toString(36).substring(2, 10);
            localStorage.setItem('chat_session_id', sessionId);
        }
        document.getElementById('sid').innerText = sessionId;

        async function startStream() {
            const prompt = document.getElementById('prompt').value;
            const outputDiv = document.getElementById('output');
            const statusDiv = document.getElementById('status');
            
            if (!prompt.trim()) return;
            outputDiv.innerText = "";
            statusDiv.innerText = "Trạng thái: Đang kết nối streaming...";

            try {
                const response = await fetch('/api/stream', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json', 'X-Session-ID': sessionId },
                    body: JSON.stringify({ prompt: prompt })
                });

                if (!response.ok) {
                    outputDiv.innerText = "Lỗi: " + await response.text();
                    statusDiv.innerText = "Trạng thái: Thất bại";
                    return;
                }

                const reader = response.body.getReader();
                const decoder = new TextDecoder();
                let count = 0;

                while (true) {
                    const { done, value } = await reader.read();
                    if (done) {
                        statusDiv.innerText = "Trạng thái: Hoàn tất (" + count + " chunks)";
                        break;
                    }
                    outputDiv.innerText += decoder.decode(value, { stream: true });
                    count++;
                }
            } catch (err) {
                outputDiv.innerText = "Lỗi mạng: " + err;
                statusDiv.innerText = "Trạng thái: Mất kết nối";
            }
        }
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
