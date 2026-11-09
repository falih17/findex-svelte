import os

class Settings:
    PROJECT_NAME: str = "Findex Biometric AI Service"
    VERSION: str = "1.0.0"
    API_V1_STR: str = "/api/v1"
    
    # Thresholds
    LIVENESS_MIN_SCORE: float = float(os.getenv("LIVENESS_MIN_SCORE", "0.40"))
    DEFAULT_SIMILARITY_THRESHOLD: float = float(os.getenv("FACE_SIMILARITY_THRESHOLD", "0.80"))
    
    # Model configs
    FACE_TARGET_SIZE: tuple = (112, 112)
    EMBEDDING_DIM: int = 512

settings = Settings()
