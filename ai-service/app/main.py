
from fastapi import FastAPI

from app.routes.diagnostic_routes import diagnostic_router
from app.routes.document_routes import document_router


app = FastAPI(
    title="NetWatch AI",
    description="AI-powered network monitoring and document intelligence service.",
    version="1.0.0",
)

app.include_router(diagnostic_router)
app.include_router(document_router)


@app.get("/api/health")
def health_check():
    return {
        "service": "NetWatch AI",
        "status": "online",
    }
