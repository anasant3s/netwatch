
import os

import faiss
import numpy as np

from app.core.config import settings


class VectorStoreService:

    def __init__(self):
        self.index = None
        self.chunk_ids = []

    def build_index(
        self,
        embeddings: list[list[float]],
        chunk_ids: list[int],
    ) -> None:
        if len(embeddings) != len(chunk_ids):
            raise ValueError(
                "Embeddings and chunk IDs must have the same length."
            )

        if not embeddings:
            self.index = None
            self.chunk_ids = []
            self.save_index()
            return

        vectors = np.asarray(embeddings, dtype="float32")

        if vectors.ndim != 2:
            raise ValueError("Embeddings must be two-dimensional.")

        if not np.isfinite(vectors).all():
            raise ValueError("Embeddings contain invalid values.")

        dimension = vectors.shape[1]

        if dimension == 0:
            raise ValueError("Embedding dimension cannot be zero.")

        faiss.normalize_L2(vectors)

        index = faiss.IndexFlatIP(dimension)
        index.add(vectors)

        self.index = index
        self.chunk_ids = chunk_ids

        self.save_index()

    def search(
        self,
        embedding: list[float],
        limit: int = 5,
    ) -> list[dict]:
        if limit <= 0:
            return []

        if self.index is None:
            self.load_index()

        if self.index is None or self.index.ntotal == 0:
            return []

        query_vector = np.asarray(
            [embedding],
            dtype="float32",
        )

        if query_vector.shape[1] != self.index.d:
            raise ValueError(
                "Query embedding dimension does not match the index."
            )

        if not np.isfinite(query_vector).all():
            raise ValueError("Query embedding contains invalid values.")

        faiss.normalize_L2(query_vector)

        scores, indices = self.index.search(
            query_vector,
            min(limit, self.index.ntotal),
        )

        results = []

        for score, index_position in zip(scores[0], indices[0]):
            if index_position < 0:
                continue

            results.append({
                "chunk_id": self.chunk_ids[index_position],
                "score": float(score),
            })

        return results

    def save_index(self) -> None:
        index_path = settings.FAISS_INDEX_PATH
        directory = os.path.dirname(index_path)

        if directory:
            os.makedirs(directory, exist_ok=True)

        metadata_path = f"{index_path}.ids"

        if self.index is None:
            if os.path.exists(index_path):
                os.remove(index_path)

            if os.path.exists(f"{metadata_path}.npy"):
                os.remove(f"{metadata_path}.npy")

            return

        faiss.write_index(self.index, index_path)

        np.save(
            metadata_path,
            np.asarray(self.chunk_ids, dtype="int64"),
        )

    def load_index(self) -> None:
        index_path = settings.FAISS_INDEX_PATH
        metadata_path = f"{index_path}.ids.npy"

        if (
            not os.path.exists(index_path)
            or not os.path.exists(metadata_path)
        ):
            return

        index = faiss.read_index(index_path)
        chunk_ids = np.load(metadata_path).tolist()

        if index.ntotal != len(chunk_ids):
            raise RuntimeError(
                "FAISS index and chunk metadata are inconsistent."
            )

        self.index = index
        self.chunk_ids = chunk_ids
