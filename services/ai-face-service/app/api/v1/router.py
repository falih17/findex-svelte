from fastapi import APIRouter
from app.api.v1.endpoints import face

api_router = APIRouter()
api_router.include_router(face.router, tags=["Face Biometrics & Liveness"])
