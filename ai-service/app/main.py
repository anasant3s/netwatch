
from fastapi import FastAPI

from app.routes.diagnostic_routes import diagnostic_router


app = FastAPI(
    title="NetWatch AI",
    description="AI-powered network incident diagnosis service.",
    version="1.0.0",
)

app.include_router(diagnostic_router)


@app.get("/api/health")
def health_check():
    return {
        "service": "NetWatch AI",
        "status": "online",
    }
