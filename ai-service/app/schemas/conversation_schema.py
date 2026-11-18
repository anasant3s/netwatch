
from datetime import datetime

from pydantic import BaseModel, Field

from app.schemas.document_schema import DocumentSource


class ConversationQuestionRequest(BaseModel):
    question: str = Field(
        min_length=3,
        max_length=1000,
    )
    conversation_id: int | None = None


class ConversationQuestionResponse(BaseModel):
    conversation_id: int
    question: str
    answer: str
    sources: list[DocumentSource]


class ConversationMessageResponse(BaseModel):
    id: int
    role: str
    content: str
    sources: list[DocumentSource] | None = None
    created_at: datetime


class ConversationResponse(BaseModel):
    id: int
    title: str
    created_at: datetime
    updated_at: datetime
    messages: list[ConversationMessageResponse]


class ConversationListItem(BaseModel):
    id: int
    title: str
    created_at: datetime
    updated_at: datetime
