
import fitz


class DocumentService:

    def extract_text(self, file_path: str) -> list[dict]:
        pages = []

        with fitz.open(file_path) as document:
            for page_number, page in enumerate(document, start=1):
                content = page.get_text("text").strip()

                if content:
                    pages.append({
                        "page_number": page_number,
                        "content": content,
                    })

        return pages
