import cv2
import numpy as np
from typing import List
from app.core.config import settings

class FeatureExtractor:
    @staticmethod
    def extract_embedding(face_img: np.ndarray) -> List[float]:
        """
        Extracts a normalized 512-dimensional embedding vector.
        In high-volume production: Uses ONNX ArcFace / MobileFaceNet session.
        """
        resized = cv2.resize(face_img, settings.FACE_TARGET_SIZE)
        gray = cv2.cvtColor(resized, cv2.COLOR_BGR2GRAY) if len(resized.shape) == 3 else resized
        
        # Spatial Grid Histogram across 4x4 patches (16 patches * 32 bins = 512 dimensions)
        hists = []
        block_h = settings.FACE_TARGET_SIZE[0] // 4
        block_w = settings.FACE_TARGET_SIZE[1] // 4
        
        for i in range(4):
            for j in range(4):
                block = gray[i*block_h:(i+1)*block_h, j*block_w:(j+1)*block_w]
                hist = cv2.calcHist([block], [0], None, [32], [0, 256]).flatten()
                hists.extend(hist)
                
        vec = np.array(hists, dtype=np.float32)
        if len(vec) > settings.EMBEDDING_DIM:
            vec = vec[:settings.EMBEDDING_DIM]
        elif len(vec) < settings.EMBEDDING_DIM:
            vec = np.pad(vec, (0, settings.EMBEDDING_DIM - len(vec)), 'constant')
            
        # L2 Normalization (Essential for Cosine Distance <=> in PostgreSQL)
        norm = np.linalg.norm(vec)
        if norm > 0:
            vec = vec / norm
            
        return [round(float(v), 6) for v in vec]
