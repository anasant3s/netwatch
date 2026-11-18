
from app.repositories.document_repository import DocumentRepository
from app.services.embedding_service import EmbeddingService
from app.services.vector_store_service import VectorStoreService


class RetrievalService:

    def __init__(self):
        self.document_repository = DocumentRepository()
        self.embedding_service = EmbeddingService()
        self.vector_store = VectorStoreService()

    def retrieve(self, query: str, limit: int = 5) -> list[dict]:
        query_embedding = self.embedding_service.generate_query_embedding(query)

        search_results = self.vector_store.search(
            embedding=query_embedding,
            limit=limit,
        )

        chunk_ids = [
            result["chunk_id"]
            for result in search_results
        ]

        chunks = self.document_repository.get_chunks_by_ids(chunk_ids)

        chunks_by_id = {
            chunk["id"]: chunk
            for chunk in chunks
        }

        results = []

        for search_result in search_results:
            chunk = chunks_by_id.get(search_result["chunk_id"])

            if chunk is None:
                continue

            results.append({
                **chunk,
                "similarity": search_result["score"],
            })

        return results
