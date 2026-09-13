package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type PromptPayload struct {
	Prompt       string  `json:"prompt"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// Gọi sang Python AI worker để lấy dữ liệu streaming
func StreamAIResponse(prompt string) (*http.Response, error) {
	// Lấy URL của service Python từ biến môi trường Docker Compose cấu hình sẵn
	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	url := fmt.Sprintf("%s/api/v1/generate", aiServiceURL)

	payload := PromptPayload{
		Prompt:      prompt,
		MaxTokens:   100,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
