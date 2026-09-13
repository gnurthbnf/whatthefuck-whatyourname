package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type AIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/stream", handleAIStream)

	log.Printf("Server đang chạy tại http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <title>Go AI Stream Pipeline</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 650px; margin: 40px auto; padding: 20px; background: #f4f4f9; }
        textarea { width: 100%; height: 100px; padding: 10px; font-size: 16px; border-radius: 5px; border: 1px solid #ccc; box-sizing: border-box; }
        button { padding: 10px 20px; font-size: 16px; background: #007bff; color: white; border: none; border-radius: 5px; cursor: pointer; margin-top: 10px; }
        button:hover { background: #0056b3; }
        #output { white-space: pre-wrap; background: white; padding: 15px; border: 1px solid #ddd; border-radius: 5px; margin-top: 20px; min-height: 120px; line-height: 1.5; }
    </style>
</head>
<body>
    <h2>Go AI Stream Pipeline (Cloud Ready)</h2>
    <textarea id="prompt" placeholder="Nhập yêu cầu của cậu...">Hãy viết một đoạn mã ngắn bằng Go để xử lý kết nối đồng thời.</textarea>
    <br>
    <button onclick="startStream()">Gửi yêu cầu Streaming</button>
    <h3>Kết quả theo thời gian thực:</h3>
    <div id="output"></div>

    <script>
        async function startStream() {
            const prompt = document.getElementById('prompt').value;
            const outputDiv = document.getElementById('output');
            outputDiv.innerText = "";

            if (!prompt.trim()) return;

            try {
                const response = await fetch('/api/stream', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ prompt: prompt })
                });

                if (!response.ok) {
                    outputDiv.innerText = "Lỗi server: " + response.statusText;
                    return;
                }

                const reader = response.body.getReader();
                const decoder = new TextDecoder();

                while (true) {
                    const { done, value } = await reader.read();
                    if (done) break;
                    outputDiv.innerText += decoder.decode(value, { stream: true });
                }
            } catch (err) {
                outputDiv.innerText = "Lỗi kết nối: " + err;
            }
        }
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func handleAIStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Chỉ hỗ trợ POST", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Dữ liệu không hợp lệ", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming không được hỗ trợ", http.StatusInternalServerError)
		return
	}

	apiKey := os.Getenv("AI_API_KEY")
	apiURL := os.Getenv("AI_API_URL")
	modelName := os.Getenv("AI_MODEL")

	if apiURL == "" {
		apiURL = "https://api.groq.com/openai/v1/chat/completions"
	}
	if modelName == "" {
		modelName = "llama-3.3-70b-versatile"
	}

	aiReq := AIRequest{
		Model: modelName,
		Messages: []Message{
			{Role: "user", Content: reqBody.Prompt},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(aiReq)
	if err != nil {
		fmt.Fprintf(w, "Lỗi tạo request AI: %v", err)
		return
	}

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(w, "Lỗi HTTP request: %v", err)
		return
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Fprintf(w, "Không thể kết nối Cloud AI: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(w, "Lỗi từ AI Provider [%d]: %s", resp.StatusCode, string(bodyBytes))
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		lineStr := string(line)
		if len(lineStr) > 5 && lineStr[:5] == "data:" {
			payload := lineStr[5:]
			if bytes.Equal(bytes.TrimSpace([]byte(payload)), []byte("[DONE]")) {
				break
			}

			var streamResp AIStreamResponse
			if err := json.Unmarshal([]byte(payload), &streamResp); err == nil {
				if len(streamResp.Choices) > 0 {
					content := streamResp.Choices[0].Delta.Content
					if content != "" {
						w.Write([]byte(content))
						flusher.Flush()
					}
				}
			}
		}
	}
}
