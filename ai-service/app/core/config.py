
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

    DATABASE_URL = os.getenv("DATABASE_URL")

    FAISS_INDEX_PATH = os.getenv(
        "FAISS_INDEX_PATH",
        "storage/faiss/netwatch.index",
    )

    GEMINI_EMBEDDING_MODEL = os.getenv(
        "GEMINI_EMBEDDING_MODEL",
        "gemini-embedding-001",
    )


settings = Settings()
