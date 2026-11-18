
import logging
import os
import tempfile

from fastapi import APIRouter, File, Form, HTTPException, UploadFile

from app.schemas.conversation_schema import (
    ConversationListItem,
    ConversationQuestionRequest,
    ConversationQuestionResponse,
    ConversationResponse,
)
from app.services.conversation_service import ConversationService
from app.services.indexing_service import IndexingService


logger = logging.getLogger(__name__)

document_router = APIRouter(
    prefix="/api/documents",
    tags=["Documents"],
)



@document_router.post("/upload")
def upload_document(
    file: UploadFile = File(...),
    reindex: bool = Form(default=False),
) -> dict:
    if not file.filename or not file.filename.lower().endswith(".pdf"):
        raise HTTPException(
            status_code=400,
            detail="Only PDF files are supported.",
        )

    temporary_path = None

    try:
        with tempfile.NamedTemporaryFile(
            delete=False,
            suffix=".pdf",
        ) as temporary_file:
            temporary_path = temporary_file.name

            file.file.seek(0)
            file_signature = file.file.read(5)

            if file_signature != b"%PDF-":
                raise HTTPException(
                    status_code=400,
                    detail="The uploaded file is not a valid PDF.",
                )

            file.file.seek(0)

            while content := file.file.read(1024 * 1024):
                temporary_file.write(content)

        indexing_service = IndexingService()

        return indexing_service.index_document(
            file_path=temporary_path,
            filename=file.filename,
            force_reindex=reindex,
        )

    except HTTPException:
        raise

    except ValueError as error:
        raise HTTPException(
            status_code=400,
            detail=str(error),
        ) from error

    except Exception as error:
        logger.exception("Failed to index uploaded document.")

        raise HTTPException(
            status_code=500,
            detail="An error occurred while processing the document.",
        ) from error

    finally:
        if temporary_path and os.path.exists(temporary_path):
            os.remove(temporary_path)

        file.file.close()



@document_router.post(
    "/ask",
    response_model=ConversationQuestionResponse,
)
def ask_document_question(
    request: ConversationQuestionRequest,
) -> ConversationQuestionResponse:
    try:
        conversation_service = ConversationService()

        result = conversation_service.ask(
            question=request.question,
            conversation_id=request.conversation_id,
        )

        return ConversationQuestionResponse.model_validate(result)

    except ValueError as error:
        raise HTTPException(
            status_code=404,
            detail=str(error),
        ) from error

    except Exception as error:
        logger.exception("Failed to answer document question.")

        raise HTTPException(
            status_code=500,
            detail="An error occurred while answering the question.",
        ) from error


@document_router.get(
    "/conversations",
    response_model=list[ConversationListItem],
)
def list_conversations() -> list[dict]:
    try:
        conversation_service = ConversationService()

        return conversation_service.get_conversations()

    except Exception as error:
        logger.exception("Failed to list conversations.")

        raise HTTPException(
            status_code=500,
            detail="An error occurred while retrieving conversations.",
        ) from error


@document_router.get(
    "/conversations/{conversation_id}",
    response_model=ConversationResponse,
)
def get_conversation(
    conversation_id: int,
) -> dict:
    try:
        conversation_service = ConversationService()

        conversation = conversation_service.get_conversation(
            conversation_id
        )

        if conversation is None:
            raise HTTPException(
                status_code=404,
                detail="Conversation not found.",
            )

        return conversation

    except HTTPException:
        raise

    except Exception as error:
        logger.exception("Failed to retrieve conversation.")

        raise HTTPException(
            status_code=500,
            detail="An error occurred while retrieving the conversation.",
        ) from error
