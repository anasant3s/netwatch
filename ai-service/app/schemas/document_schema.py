
from pydantic import BaseModel, Field


class DocumentQuestionRequest(BaseModel):
    question: str = Field(
        min_length=3,
        max_length=1000,
        description="Question about the indexed technical documents.",
    )


class DocumentSource(BaseModel):
    document_id: int
    filename: str
    page_number: int
    similarity: float


class DocumentQuestionResponse(BaseModel):
    question: str
    answer: str
    sources: list[DocumentSource]
