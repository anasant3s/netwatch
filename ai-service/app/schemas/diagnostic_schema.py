
from pydantic import BaseModel, Field


class IncidentRequest(BaseModel):
    device_id: int
    device_name: str
    ip: str
    previous_status: str
    current_status: str
    message: str
    created_at: str


class DiagnosticResponse(BaseModel):
    device_id: int
    device_name: str
    diagnosis: str
    possible_causes: list[str] = Field(default_factory=list)
    recommended_actions: list[str] = Field(default_factory=list)
    confidence: str
