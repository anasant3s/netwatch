
from google import genai
from google.genai import types

from app.core.config import settings


class EmbeddingService:

    def generate_embedding(self, text: str) -> list[float]:
        if not settings.GEMINI_API_KEY:
            raise RuntimeError("GEMINI_API_KEY is not configured.")

        with genai.Client(api_key=settings.GEMINI_API_KEY) as client:
            response = client.models.embed_content(
                model=settings.GEMINI_EMBEDDING_MODEL,
                contents=text,
                config=types.EmbedContentConfig(
                    task_type="RETRIEVAL_DOCUMENT",
                    output_dimensionality=768,
                ),
            )

        if not response.embeddings:
            raise RuntimeError("Failed to generate text embedding.")

        return response.embeddings[0].values

    def generate_query_embedding(self, text: str) -> list[float]:
        if not settings.GEMINI_API_KEY:
            raise RuntimeError("GEMINI_API_KEY is not configured.")

        with genai.Client(api_key=settings.GEMINI_API_KEY) as client:
            response = client.models.embed_content(
                model=settings.GEMINI_EMBEDDING_MODEL,
                contents=text,
                config=types.EmbedContentConfig(
                    task_type="RETRIEVAL_QUERY",
                    output_dimensionality=768,
                ),
            )

        if not response.embeddings:
            raise RuntimeError("Failed to generate query embedding.")

        return response.embeddings[0].values
