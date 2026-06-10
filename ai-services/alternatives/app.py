import os
import time
import json
import random
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
from fastapi.middleware.cors import CORSMiddleware
import redis

app = FastAPI(title="Yanwit AI - Alternatives Generator")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://127.0.0.1:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Подключение к Redis для кэширования
redis_client = redis.Redis(
    host=os.getenv("REDIS_HOST", "localhost"),
    port=int(os.getenv("REDIS_PORT", 6379)),
    decode_responses=True
)

class AlternativeRequest(BaseModel):
    text: str
    count: int = 3
    style: str = "funny"  # funny, professional, sarcastic, encouraging

class AlternativeResponse(BaseModel):
    alternatives: List[str]
    from_cache: bool = False
    latency_ms: int = 0

# Шаблоны для разных стилей (заглушка, позже заменим на реальную ML модель)
STYLE_TEMPLATES = {
    "funny": [
        "Знаете, {text}... 😅",
        "{text} 🚀 (шутка)",
        "Слушайте, {text} 🔥"
    ],
    "professional": [
        "Важно отметить: {text}",
        "Следует подчеркнуть: {text}",
        "Обратите внимание: {text}"
    ],
    "sarcastic": [
        "О да, {text}, конечно 😏",
        "{text}... ну-ну 🤔",
        "Очевидно же, что {text} 🙄"
    ],
    "encouraging": [
        "Вы правы! {text} 💪",
        "Отличная мысль! {text} 🌟",
        "Поддерживаю! {text} 🎯"
    ]
}

@app.get("/health")
async def health():
    return {"status": "ok", "service": "alternatives"}

@app.post("/generate", response_model=AlternativeResponse)
async def generate_alternatives(request: AlternativeRequest):
    start_time = time.time()
    
    # Проверяем кэш
    cache_key = f"alt:{hash(request.text)}:{request.style}"
    cached = redis_client.get(cache_key)
    
    if cached:
        return AlternativeResponse(
            alternatives=json.loads(cached),
            from_cache=True,
            latency_ms=int((time.time() - start_time) * 1000)
        )
    
    # Генерация альтернатив (заглушка)
    templates = STYLE_TEMPLATES.get(request.style, STYLE_TEMPLATES["funny"])
    alternatives = []
    
    for i in range(request.count):
        template = templates[i % len(templates)]
        alt = template.format(text=request.text[:50])
        alternatives.append(alt)
    
    # Кэшируем результат на 1 час
    redis_client.setex(cache_key, 3600, json.dumps(alternatives))
    
    return AlternativeResponse(
        alternatives=alternatives,
        from_cache=False,
        latency_ms=int((time.time() - start_time) * 1000)
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8002)