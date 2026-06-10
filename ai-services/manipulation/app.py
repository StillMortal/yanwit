import os
import time
import json
import re
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Optional
from fastapi.middleware.cors import CORSMiddleware
import redis

app = FastAPI(title="Yanwit AI - Manipulation Detector")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://127.0.0.1:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

redis_client = redis.Redis(
    host=os.getenv("REDIS_HOST", "localhost"),
    port=int(os.getenv("REDIS_PORT", 6379)),
    decode_responses=True
)

class ManipulationRequest(BaseModel):
    text: str

class ManipulationResponse(BaseModel):
    has_manipulation: bool
    types: List[str]
    confidence: float
    suggestions: List[str]
    rewritten_version: Optional[str] = None

# Ключевые слова для детекции манипуляций (заглушка)
MANIPULATION_PATTERNS = {
    "bandwagon": {
        "keywords": ["все так делают", "все нормальные люди", "каждый", "миллионы", "никто не", "всем нравится", "всем нравятся"],
        "suggestion": "Используйте 'мне нравится' вместо 'всем нравятся'"
    },
    "authority": {
        "keywords": ["эксперты говорят", "ученые доказали", "по данным", "исследования показывают"],
        "suggestion": "Приведите конкретные данные вместо ссылки на авторитет"
    },
    "fear": {
        "keywords": ["иначе будет", "если не сделаешь", "пожалеешь", "упустишь возможность"],
        "suggestion": "Опишите позитивные последствия вместо негативных"
    },
    "urgency": {
        "keywords": ["срочно", "только сегодня", "последний шанс", "ограниченное предложение"],
        "suggestion": "Уберите искусственное ограничение по времени"
    },
    "scarcity": {
        "keywords": ["осталось всего", "спешите", "дефицит", "ограниченное количество"],
        "suggestion": "Объясните реальную, а не искусственную ценность"
    }
}

@app.get("/health")
async def health():
    return {"status": "ok", "service": "manipulation"}

@app.post("/detect", response_model=ManipulationResponse)
async def detect_manipulation(request: ManipulationRequest):
    start_time = time.time()
    text_lower = request.text.lower()
    
    # Проверяем кэш
    cache_key = f"manip:{hash(request.text)}"
    cached = redis_client.get(cache_key)
    
    if cached:
        result = json.loads(cached)
        return ManipulationResponse(**result)
    
    detected_types = []
    suggestions = []
    confidence = 0.0
    
    for manip_type, pattern in MANIPULATION_PATTERNS.items():
        for keyword in pattern["keywords"]:
            if keyword in text_lower:
                detected_types.append(manip_type)
                suggestions.append(pattern["suggestion"])
                confidence += 0.25
                break
    
    confidence = min(confidence, 0.95)
    has_manipulation = len(detected_types) > 0
    
    result = {
        "has_manipulation": has_manipulation,
        "types": detected_types,
        "confidence": confidence,
        "suggestions": suggestions,
        "rewritten_version": None
    }
    
    # Кэшируем результат на 6 часов
    redis_client.setex(cache_key, 21600, json.dumps(result))
    
    return ManipulationResponse(**result)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8003)