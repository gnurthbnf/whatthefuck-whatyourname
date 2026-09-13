import os

class Settings:
    PROJECT_NAME: str = "AI Worker Service"
    HOST: str = os.getenv("HOST", "0.0.0.0")
    PORT: int = int(os.getenv("PORT", 8000))
    MODEL_NAME: str = os.getenv("MODEL_NAME", "default-stream-model")

settings = Settings()
