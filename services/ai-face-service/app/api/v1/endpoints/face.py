import time
import cv2
import numpy as np
from fastapi import APIRouter, File, UploadFile, HTTPException, status
from app.core.config import settings
from app.schemas.face import FaceExtractResponse, VectorCompareRequest, VectorCompareResponse
from app.models.detector import FaceDetector
from app.models.liveness import LivenessDetector
from app.models.extractor import FeatureExtractor

router = APIRouter()
detector = FaceDetector()

@router.post("/extract-face", response_model=FaceExtractResponse)
async def extract_face(file: UploadFile = File(...)):
    start_time = time.time()
    contents = await file.read()
    
    if len(contents) == 0:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="File gambar kosong")
        
    nparr = np.frombuffer(contents, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    if img is None:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Format file bukan gambar yang valid")
        
    detection_res = detector.detect_largest_face(img)
    if detection_res is None:
        elapsed = (time.time() - start_time) * 1000
        return FaceExtractResponse(
            success=False,
            face_detected=False,
            liveness_score=0.0,
            is_live=False,
            embedding=[],
            processing_time_ms=round(elapsed, 2),
            message="Wajah tidak ditemukan pada gambar."
        )
        
    face_crop, _ = detection_res
    liveness_score = LivenessDetector.evaluate(face_crop)
    is_live = liveness_score >= settings.LIVENESS_MIN_SCORE
    embedding = FeatureExtractor.extract_embedding(face_crop)
    
    elapsed = (time.time() - start_time) * 1000
    return FaceExtractResponse(
        success=True,
        face_detected=True,
        liveness_score=liveness_score,
        is_live=is_live,
        embedding=embedding,
        processing_time_ms=round(elapsed, 2),
        message="Fitur wajah dan liveness berhasil diekstraksi."
    )

@router.post("/compare-vectors", response_model=VectorCompareResponse)
def compare_vectors(req: VectorCompareRequest, threshold: float = settings.DEFAULT_SIMILARITY_THRESHOLD):
    if len(req.vector_a) != settings.EMBEDDING_DIM or len(req.vector_b) != settings.EMBEDDING_DIM:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, 
            detail=f"Vector harus tepat berdimensi {settings.EMBEDDING_DIM}"
        )
        
    a = np.array(req.vector_a, dtype=np.float32)
    b = np.array(req.vector_b, dtype=np.float32)
    
    dot_prod = np.dot(a, b)
    norm_a = np.linalg.norm(a)
    norm_b = np.linalg.norm(b)
    
    sim = dot_prod / (norm_a * norm_b) if (norm_a > 0 and norm_b > 0) else 0.0
    sim = float(np.clip(sim, 0.0, 1.0))
    
    return VectorCompareResponse(
        similarity_score=round(sim, 4),
        is_match=sim >= threshold,
        threshold_used=threshold
    )
