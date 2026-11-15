
import os

from dotenv import load_dotenv


load_dotenv()


class Settings:
    GEMINI_API_KEY = os.getenv("GEMINI_API_KEY")
    GEMINI_MODEL = os.getenv("GEMINI_MODEL", "gemini-3.8-flash")
    NETWATCH_API_URL = os.getenv(
        "NETWATCH_API_URL",
        "http://localhost:8080",
    )


settings = Settings()
