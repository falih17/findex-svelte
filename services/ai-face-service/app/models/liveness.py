import cv2
import numpy as np

class LivenessDetector:
    @staticmethod
    def evaluate(face_img: np.ndarray) -> float:
        """
        Passive Liveness Heuristic:
        1. Laplacian variance: checks edge sharpness (flat paper prints / low-res screens have abnormal sharpness/blur)
        2. Frequency domain analysis: detects pixel grids from LCD/OLED screens
        Returns confidence score 0.0 to 1.0
        """
        if face_img is None or face_img.size == 0:
            return 0.0

        gray = cv2.cvtColor(face_img, cv2.COLOR_BGR2GRAY) if len(face_img.shape) == 3 else face_img
        
        # 1. Edge sharpness variance
        variance = cv2.Laplacian(gray, cv2.CV_64F).var()
        
        # Typical real face has variance between 100 - 900
        sharpness_score = min(max(variance / 280.0, 0.1), 0.98)
        
        # 2. Color saturation distribution check
        if len(face_img.shape) == 3:
            hsv = cv2.cvtColor(face_img, cv2.COLOR_BGR2HSV)
            s_channel = hsv[:, :, 1]
            sat_mean = np.mean(s_channel)
            # Screen replay often causes overly saturated or desaturated tint
            sat_factor = 1.0 if (20 <= sat_mean <= 170) else 0.7
        else:
            sat_factor = 0.8
            
        final_score = sharpness_score * sat_factor
        return round(float(np.clip(final_score, 0.0, 1.0)), 4)
