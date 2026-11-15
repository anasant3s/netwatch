import time
from google import genai
from google.genai import types

from google.genai import errors
from app.core.config import settings
from app.schemas.diagnostic_schema import (
    IncidentRequest,
    DiagnosticResponse,
)


class GeminiService:
    def __init__(self):
        if not settings.GEMINI_API_KEY:
            raise ValueError("GEMINI_API_KEY is not configured.")

    
    def diagnose(self, incident: IncidentRequest) -> DiagnosticResponse:
        prompt = f"""
        Você é um assistente especializado em diagnóstico
        de infraestrutura de redes.

        Analise o incidente:

        Equipamento: {incident.device_name}
        ID: {incident.device_id}
        IP: {incident.ip}
        Status anterior: {incident.previous_status}
        Status atual: {incident.current_status}
        Mensagem: {incident.message}
        Data: {incident.created_at}

        Instruções:
        - Responda em português brasileiro.
        - Analise somente as informações fornecidas.
        - Não apresente hipóteses como fatos.
        - Recomende ações práticas de diagnóstico.
        - Informe quando faltarem evidências.
        - A confiança deve ser low, medium ou high.
        """

        for attempt in range(3):
            try:
                with genai.Client(
                    api_key=settings.GEMINI_API_KEY
                ) as client:
                    response = client.models.generate_content(
                        model=settings.GEMINI_MODEL,
                        contents=prompt,
                        config=types.GenerateContentConfig(
                            response_mime_type="application/json",
                            response_schema=DiagnosticResponse,
                            temperature=0.2,
                        ),
                    )

                if response.parsed is not None:
                    return DiagnosticResponse.model_validate(
                        response.parsed
                    )

                if not response.text:
                    raise RuntimeError(
                        "Gemini returned an empty response."
                    )

                return DiagnosticResponse.model_validate_json(
                    response.text
                )

            except errors.ServerError as error:
                if attempt == 2:
                    raise

                time.sleep(2 ** attempt)

        raise RuntimeError("Gemini diagnosis failed.")
