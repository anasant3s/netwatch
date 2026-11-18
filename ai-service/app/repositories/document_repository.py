
from psycopg.types.json import Jsonb

from app.core.database import get_connection


class DocumentRepository:

    def create_document(self, filename: str, file_hash: str) -> int:
        query = """
            INSERT INTO documents (filename, file_hash)
            VALUES (%s, %s)
            RETURNING id;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (filename, file_hash))
                return cursor.fetchone()[0]

    def get_document_by_hash(self, file_hash: str) -> dict | None:
        query = """
            SELECT id, filename, file_hash, status, created_at
            FROM documents
            WHERE file_hash = %s;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (file_hash,))
                row = cursor.fetchone()

        if row is None:
            return None

        return {
            "id": row[0],
            "filename": row[1],
            "file_hash": row[2],
            "status": row[3],
            "created_at": row[4],
        }

    def prepare_document_for_indexing(self, document_id: int) -> None:
        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(
                    "DELETE FROM document_chunks WHERE document_id = %s;",
                    (document_id,),
                )

                cursor.execute(
                    """
                    UPDATE documents
                    SET status = 'pending'
                    WHERE id = %s;
                    """,
                    (document_id,),
                )

    def update_document_status(
        self,
        document_id: int,
        status: str,
    ) -> None:
        query = """
            UPDATE documents
            SET status = %s
            WHERE id = %s;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (status, document_id))

    def create_chunks(self, document_id: int, chunks: list[dict]) -> None:
        query = """
            INSERT INTO document_chunks (
                document_id,
                chunk_index,
                page_number,
                content,
                embedding
            )
            VALUES (%s, %s, %s, %s, %s);
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                for chunk in chunks:
                    embedding = chunk.get("embedding")

                    cursor.execute(
                        query,
                        (
                            document_id,
                            chunk["chunk_index"],
                            chunk["page_number"],
                            chunk["content"],
                            Jsonb(embedding) if embedding is not None else None,
                        ),
                    )

    def replace_chunks(
        self,
        document_id: int,
        chunks: list[dict],
    ) -> None:
        delete_query = """
            DELETE FROM document_chunks
            WHERE document_id = %s;
        """

        insert_query = """
            INSERT INTO document_chunks (
                document_id,
                chunk_index,
                page_number,
                content,
                embedding
            )
            VALUES (%s, %s, %s, %s, %s);
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(delete_query, (document_id,))

                for chunk in chunks:
                    cursor.execute(
                        insert_query,
                        (
                            document_id,
                            chunk["chunk_index"],
                            chunk["page_number"],
                            chunk["content"],
                            Jsonb(chunk["embedding"]),
                        ),
                    )

                cursor.execute(
                    """
                    UPDATE documents
                    SET status = 'indexed'
                    WHERE id = %s;
                    """,
                    (document_id,),
                )

    def get_document(self, document_id: int) -> dict | None:
        query = """
            SELECT id, filename, file_hash, status, created_at
            FROM documents
            WHERE id = %s;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (document_id,))
                row = cursor.fetchone()

        if row is None:
            return None

        return {
            "id": row[0],
            "filename": row[1],
            "file_hash": row[2],
            "status": row[3],
            "created_at": row[4],
        }

    def get_chunks_by_ids(self, chunk_ids: list[int]) -> list[dict]:
        if not chunk_ids:
            return []

        query = """
            SELECT
                chunks.id,
                chunks.document_id,
                chunks.page_number,
                chunks.content
            FROM document_chunks AS chunks
            INNER JOIN documents AS documents
                ON documents.id = chunks.document_id
            WHERE chunks.id = ANY(%s)
              AND documents.status = 'indexed';
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (chunk_ids,))
                rows = cursor.fetchall()

        chunks_by_id = {
            row[0]: {
                "id": row[0],
                "document_id": row[1],
                "page_number": row[2],
                "content": row[3],
            }
            for row in rows
        }

        return [
            chunks_by_id[chunk_id]
            for chunk_id in chunk_ids
            if chunk_id in chunks_by_id
        ]

    def get_all_embeddings(self) -> list[dict]:
        query = """
            SELECT
                chunks.id,
                chunks.embedding
            FROM document_chunks AS chunks
            INNER JOIN documents AS documents
                ON documents.id = chunks.document_id
            WHERE chunks.embedding IS NOT NULL
              AND documents.status = 'indexed'
            ORDER BY chunks.id;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query)
                rows = cursor.fetchall()

        return [
            {
                "chunk_id": row[0],
                "embedding": row[1],
            }
            for row in rows
        ]
