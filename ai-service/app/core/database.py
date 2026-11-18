
import psycopg

from app.core.config import settings


def get_connection():
    if not settings.DATABASE_URL:
        raise RuntimeError("DATABASE_URL is not configured.")

    return psycopg.connect(settings.DATABASE_URL)


def check_database_connection() -> bool:
    with get_connection() as connection:
        with connection.cursor() as cursor:
            cursor.execute("SELECT 1")
            return cursor.fetchone()[0] == 1
