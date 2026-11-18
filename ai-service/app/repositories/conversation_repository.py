
from psycopg.types.json import Jsonb

from app.core.database import get_connection


class ConversationRepository:

    def create_conversation(self, title: str) -> int:
        query = """
            INSERT INTO conversations (title)
            VALUES (%s)
            RETURNING id;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (title,))
                return cursor.fetchone()[0]

    def get_conversation(self, conversation_id: int) -> dict | None:
        query = """
            SELECT id, title, created_at, updated_at
            FROM conversations
            WHERE id = %s;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (conversation_id,))
                row = cursor.fetchone()

        if row is None:
            return None

        return {
            "id": row[0],
            "title": row[1],
            "created_at": row[2],
            "updated_at": row[3],
        }

    def get_conversations(self) -> list[dict]:
        query = """
            SELECT id, title, created_at, updated_at
            FROM conversations
            ORDER BY updated_at DESC;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query)
                rows = cursor.fetchall()

        return [
            {
                "id": row[0],
                "title": row[1],
                "created_at": row[2],
                "updated_at": row[3],
            }
            for row in rows
        ]

    def create_message(
        self,
        conversation_id: int,
        role: str,
        content: str,
        sources: list[dict] | None = None,
    ) -> int:
        query = """
            INSERT INTO messages (
                conversation_id,
                role,
                content,
                sources
            )
            VALUES (%s, %s, %s, %s)
            RETURNING id;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(
                    query,
                    (
                        conversation_id,
                        role,
                        content,
                        Jsonb(sources) if sources is not None else None,
                    ),
                )

                return cursor.fetchone()[0]

    def get_messages(self, conversation_id: int) -> list[dict]:
        query = """
            SELECT id, role, content, sources, created_at
            FROM messages
            WHERE conversation_id = %s
            ORDER BY created_at ASC, id ASC;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (conversation_id,))
                rows = cursor.fetchall()

        return [
            {
                "id": row[0],
                "role": row[1],
                "content": row[2],
                "sources": row[3],
                "created_at": row[4],
            }
            for row in rows
        ]

    def update_conversation_timestamp(
        self,
        conversation_id: int,
    ) -> None:
        query = """
            UPDATE conversations
            SET updated_at = CURRENT_TIMESTAMP
            WHERE id = %s;
        """

        with get_connection() as connection:
            with connection.cursor() as cursor:
                cursor.execute(query, (conversation_id,))
