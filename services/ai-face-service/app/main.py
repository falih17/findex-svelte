import time
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.core.config import settings
from app.api.v1.router import api_router

app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
    description="Findex Biometric Microservice for Anti-Spoofing & 512-D Feature Vector Extraction"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
def health_check():
    return {
        "app": settings.PROJECT_NAME,
        "status": "healthy",
        "timestamp": time.time(),
        "version": settings.VERSION
    }

app.include_router(api_router, prefix=settings.API_V1_STR)
