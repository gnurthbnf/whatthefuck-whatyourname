package client

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type ServiceClient struct {
	httpClient   *http.Client
	aiServiceURL string
	rustURL      string
}

func NewServiceClient() *ServiceClient {
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://localhost:8000"
	}

	rustURL := os.Getenv("RUST_SERVICE_URL")
	if rustURL == "" {
		rustURL = "http://localhost:8001"
	}

	return &ServiceClient{
		httpClient:   &http.Client{Timeout: 5 * time.Second},
		aiServiceURL: aiURL,
		rustURL:      rustURL,
	}
}

func (c *ServiceClient) CallAIWorker() error {
	resp, err := c.httpClient.Get(c.aiServiceURL + "/health")
	if err != nil {
		return fmt.Errorf("không thể kết nối AI worker: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func (c *ServiceClient) CallRustWorker() error {
	resp, err := c.httpClient.Get(c.rustURL + "/health")
	if err != nil {
		return fmt.Errorf("không thể kết nối Rust worker: %w", err)
	}
	defer resp.Body.Close()
	return nil
}
