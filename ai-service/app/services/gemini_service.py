
import time

from google import genai
from google.genai import errors, types

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

            except errors.ServerError:
                if attempt == 2:
                    raise

                time.sleep(2 ** attempt)

        raise RuntimeError("Gemini diagnosis failed.")

    def answer_question(
        self,
        question: str,
        context_chunks: list[dict],
        conversation_history: list[dict] | None = None,
    ) -> str:
        if not context_chunks:
            return (
                "Não encontrei informações relevantes "
                "nos documentos disponíveis."
            )

        context = "\n\n".join(
            (
                f"Documento ID: {chunk['document_id']}\n"
                f"Página: {chunk['page_number']}\n"
                f"Conteúdo:\n{chunk['content']}"
            )
            for chunk in context_chunks
        )

        history = ""

        if conversation_history:
            history = "\n\n".join(
                (
                    f"{'Usuário' if message['role'] == 'user' else 'Assistente'}: "
                    f"{message['content']}"
                )
                for message in conversation_history[-10:]
            )

        prompt = f"""
        Você é um assistente especializado em documentação
        técnica de infraestrutura de redes e segurança eletrônica.

        Responda à pergunta utilizando exclusivamente as informações
        técnicas presentes no contexto dos documentos.

        HISTÓRICO DA CONVERSA:
        {history if history else "Esta é a primeira pergunta da conversa."}

        PERGUNTA ATUAL:
        {question}

        CONTEXTO DOS DOCUMENTOS:
        {context}

        INSTRUÇÕES:
        - Responda em português brasileiro.
        - Utilize somente informações técnicas presentes
          nos documentos fornecidos.
        - Utilize o histórico para compreender referências,
          complementações e perguntas relacionadas.
        - Não trate respostas anteriores como evidências técnicas.
        - Não invente especificações ou características técnicas.
        - Se o contexto não responder à pergunta, informe isso.
        - Seja objetivo e tecnicamente preciso.
        - Sempre que possível, indique o documento e a página
          que fundamentam a resposta.
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
                            temperature=0.2,
                        ),
                    )

                if not response.text:
                    raise RuntimeError(
                        "Gemini returned an empty response."
                    )

                return response.text.strip()

            except errors.ServerError:
                if attempt == 2:
                    raise

                time.sleep(2 ** attempt)

        raise RuntimeError("Document question answering failed.")
