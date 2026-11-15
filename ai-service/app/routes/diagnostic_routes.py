
import logging

from fastapi import APIRouter, HTTPException

from app.schemas.diagnostic_schema import (
    IncidentRequest,
    DiagnosticResponse,
)
from app.services.gemini_service import GeminiService
from app.services.netwatch_service import (
    NetWatchAPIError,
    NetWatchService,
)


logger = logging.getLogger(__name__)

diagnostic_router = APIRouter(
    prefix="/api/diagnostics",
    tags=["Diagnostics"],
)


@diagnostic_router.post(
    "",
    response_model=DiagnosticResponse,
)
def create_diagnostic(incident: IncidentRequest):
    try:
        service = GeminiService()
        return service.diagnose(incident)

    except ValueError as error:
        logger.exception("Invalid diagnostic configuration.")
        raise HTTPException(
            status_code=500,
            detail=str(error),
        )

    except Exception:
        logger.exception("Failed to generate incident diagnosis.")
        raise HTTPException(
            status_code=502,
            detail="Failed to generate incident diagnosis.",
        )


@diagnostic_router.get(
    "/devices/{device_id}/incidents/{incident_id}",
    response_model=DiagnosticResponse,
)
def diagnose_registered_incident(
    device_id: int,
    incident_id: int,
):
    try:
        netwatch_service = NetWatchService()

        incident = netwatch_service.get_incident(
            device_id=device_id,
            incident_id=incident_id,
        )

        gemini_service = GeminiService()

        return gemini_service.diagnose(incident)

    except NetWatchAPIError as error:
        raise HTTPException(
            status_code=error.status_code,
            detail=str(error),
        )

    except ValueError as error:
        logger.exception("Invalid diagnostic configuration.")
        raise HTTPException(
            status_code=500,
            detail=str(error),
        )

    except Exception:
        logger.exception("Failed to diagnose registered incident.")
        raise HTTPException(
            status_code=502,
            detail="Failed to generate incident diagnosis.",
        )
