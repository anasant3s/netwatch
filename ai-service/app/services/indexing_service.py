
import hashlib

from app.repositories.document_repository import DocumentRepository
from app.services.chunk_service import ChunkService
from app.services.document_service import DocumentService
from app.services.embedding_service import EmbeddingService
from app.services.vector_store_service import VectorStoreService


class IndexingService:

    def __init__(self):
        self.document_repository = DocumentRepository()
        self.document_service = DocumentService()
        self.chunk_service = ChunkService()
        self.embedding_service = EmbeddingService()
        self.vector_store = VectorStoreService()

    def index_document(
        self,
        file_path: str,
        filename: str,
        force_reindex: bool = False,
    ) -> dict:
        file_hash = self._calculate_file_hash(file_path)

        existing_document = (
            self.document_repository.get_document_by_hash(file_hash)
        )

        if existing_document:
            document_id = existing_document["id"]

            if (
                existing_document["status"] == "indexed"
                and not force_reindex
            ):
                return {
                    "document_id": document_id,
                    "filename": existing_document["filename"],
                    "status": "already_indexed",
                }

        else:
            document_id = self.document_repository.create_document(
                filename=filename,
                file_hash=file_hash,
            )

        try:
            pages = self.document_service.extract_text(file_path)

            if not pages:
                raise ValueError("No readable text found in the PDF.")

            chunks = self.chunk_service.split_pages(pages)

            if not chunks:
                raise ValueError("No text chunks generated from the PDF.")

            for chunk in chunks:
                chunk["embedding"] = (
                    self.embedding_service.generate_embedding(
                        chunk["content"]
                    )
                )

            self.document_repository.replace_chunks(
                document_id,
                chunks,
            )

            self._rebuild_vector_index()

            return {
                "document_id": document_id,
                "filename": filename,
                "total_pages": len(pages),
                "total_chunks": len(chunks),
                "status": "indexed",
            }

        except Exception:
            self.document_repository.update_document_status(
                document_id,
                "failed",
            )
            raise

    def _rebuild_vector_index(self) -> None:
        stored_chunks = self.document_repository.get_all_embeddings()

        embeddings = [
            chunk["embedding"]
            for chunk in stored_chunks
        ]

        chunk_ids = [
            chunk["chunk_id"]
            for chunk in stored_chunks
        ]

        self.vector_store.build_index(embeddings, chunk_ids)

    @staticmethod
    def _calculate_file_hash(file_path: str) -> str:
        file_digest = hashlib.sha256()

        with open(file_path, "rb") as file:
            for block in iter(lambda: file.read(8192), b""):
                file_digest.update(block)

        return file_digest.hexdigest()
