
class ChunkService:

    def __init__(
        self,
        chunk_size: int = 1000,
        chunk_overlap: int = 150,
    ):
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

    def split_pages(self, pages: list[dict]) -> list[dict]:
        chunks = []
        chunk_index = 0

        for page in pages:
            content = page["content"]
            start = 0

            while start < len(content):
                end = start + self.chunk_size
                chunk_content = content[start:end].strip()

                if chunk_content:
                    chunks.append({
                        "chunk_index": chunk_index,
                        "page_number": page["page_number"],
                        "content": chunk_content,
                    })
                    chunk_index += 1

                start += self.chunk_size - self.chunk_overlap

        return chunks
