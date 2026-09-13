import unittest
import httpx
import os

class TestAIWorker(unittest.TestCase):
    BASE_URL = os.getenv("AI_SERVICE_URL", "http://localhost:8000")

    def test_health_endpoint(self):
        """Kiểm tra endpoint health check của AI worker"""
        response = httpx.get(f"{self.BASE_URL}/health")
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertEqual(data["status"], "healthy")

    def test_generation_stream(self):
        """Kiểm tra tính năng streaming phản hồi AI"""
        payload = {"prompt": "Xin chào hệ thống", "max_tokens": 50, "temperature": 0.5}
        with httpx.stream("POST", f"{self.BASE_URL}/api/v1/generate", json=payload, timeout=10.0) as response:
            self.assertEqual(response.status_code, 200)
            chunks = list(response.iter_text())
            self.assertTrue(len(chunks) > 0)

if __name__ == "__main__":
    unittest.main()
