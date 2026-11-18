
from app.repositories.conversation_repository import ConversationRepository
from app.repositories.document_repository import DocumentRepository
from app.services.gemini_service import GeminiService
from app.services.retrieval_service import RetrievalService


class ConversationService:

    def __init__(self):
        self.conversation_repository = ConversationRepository()
        self.document_repository = DocumentRepository()
        self.retrieval_service = RetrievalService()
        self.gemini_service = GeminiService()

    def ask(self, question: str, conversation_id: int | None = None) -> dict:
        if conversation_id is None:
            conversation_id = (
                self.conversation_repository.create_conversation(
                    title=self._generate_title(question)
                )
            )

            conversation_history = []

        else:
            conversation = (
                self.conversation_repository.get_conversation(
                    conversation_id
                )
            )

            if conversation is None:
                raise ValueError("Conversation not found.")

            conversation_history = (
                self.conversation_repository.get_messages(
                    conversation_id
                )
            )

        self.conversation_repository.create_message(
            conversation_id=conversation_id,
            role="user",
            content=question,
        )

        context_chunks = self.retrieval_service.retrieve(
            query=question,
            limit=5,
        )

        answer = self.gemini_service.answer_question(
            question=question,
            context_chunks=context_chunks,
            conversation_history=conversation_history,
        )

        sources = self._build_sources(context_chunks)

        self.conversation_repository.create_message(
            conversation_id=conversation_id,
            role="assistant",
            content=answer,
            sources=sources,
        )

        self.conversation_repository.update_conversation_timestamp(
            conversation_id
        )

        return {
            "conversation_id": conversation_id,
            "question": question,
            "answer": answer,
            "sources": sources,
        }

    def _build_sources(self, context_chunks: list[dict]) -> list[dict]:
        sources = []
        seen_pages = set()

        for chunk in context_chunks:
            document_id = chunk["document_id"]
            page_number = chunk["page_number"]

            page_key = (document_id, page_number)

            if page_key in seen_pages:
                continue

            document = self.document_repository.get_document(
                document_id
            )

            if document is None:
                continue

            seen_pages.add(page_key)

            sources.append({
                "document_id": document_id,
                "filename": document["filename"],
                "page_number": page_number,
                "similarity": chunk["similarity"],
            })

        return sources

    def get_conversations(self) -> list[dict]:
        return self.conversation_repository.get_conversations()

    def get_conversation(self, conversation_id: int) -> dict | None:
        conversation = (
            self.conversation_repository.get_conversation(
                conversation_id
            )
        )

        if conversation is None:
            return None

        messages = self.conversation_repository.get_messages(
            conversation_id
        )

        conversation["messages"] = messages

        return conversation

    @staticmethod
    def _generate_title(question: str) -> str:
        title = question.strip()

        if len(title) > 80:
            return title[:77] + "..."

        return title
