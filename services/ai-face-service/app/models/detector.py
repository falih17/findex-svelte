import cv2
import numpy as np
from typing import Optional, Tuple

class FaceDetector:
    def __init__(self):
        cascade_path = cv2.data.haarcascades + 'haarcascade_frontalface_default.xml'
        self.cascade = cv2.CascadeClassifier(cascade_path)

    def detect_largest_face(self, img: np.ndarray) -> Optional[Tuple[np.ndarray, Tuple[int, int, int, int]]]:
        if img is None or img.size == 0:
            return None
            
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY) if len(img.shape) == 3 else img
        faces = self.cascade.detectMultiScale(
            gray, 
            scaleFactor=1.1, 
            minNeighbors=5, 
            minSize=(60, 60)
        )
        
        if len(faces) == 0:
            return None
            
        # Sort by bounding box area (largest face in foreground)
        faces = sorted(faces, key=lambda f: f[2] * f[3], reverse=True)
        x, y, w, h = faces[0]
        
        # Add slight margin padding around the face
        pad_y = int(h * 0.1)
        pad_x = int(w * 0.1)
        
        y1 = max(0, y - pad_y)
        y2 = min(img.shape[0], y + h + pad_y)
        x1 = max(0, x - pad_x)
        x2 = min(img.shape[1], x + w + pad_x)
        
        face_crop = img[y1:y2, x1:x2]
        return face_crop, (x, y, w, h)
