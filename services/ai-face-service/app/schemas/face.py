from pydantic import BaseModel, Field
from typing import List

class FaceExtractResponse(BaseModel):
    success: bool
    face_detected: bool
    liveness_score: float = Field(..., description="Calculated passive liveness score between 0.0 and 1.0")
    is_live: bool = Field(..., description="True if liveness score exceeds anti-spoofing threshold")
    embedding: List[float] = Field(default=[], description="512-dimensional vector embedding")
    processing_time_ms: float
    message: str

class VectorCompareRequest(BaseModel):
    vector_a: List[float]
    vector_b: List[float]

class VectorCompareResponse(BaseModel):
    similarity_score: float
    is_match: bool
    threshold_used: float
