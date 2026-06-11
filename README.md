# 🐦 Yanwit — микроблог с AI-суперсилами

**Yanwit** — это современное приложение для ведения микроблога с интегрированными AI-функциями. Проект создан для демонстрации возможностей Go, Python, Docker и современных веб-технологий.

## ✨ Особенности

### 🚀 Базовый функционал
- Регистрация и JWT-аутентификация
- Публикация твитов (до 280 символов)
- Лента подписок (home timeline)
- Подписка и отписка от пользователей
- Лайки и счётчики
- Поиск пользователей
- Тёмная и светлая темы (сохраняется в localStorage)

### 🤖 AI-функции (уникальные)

#### 1. Генератор альтернативных концовок ✨
Напишите черновик твита, нажмите **"✨ AI улучшить"** — и получите 3 стилистически разных варианта:
- 😄 **Весёлый** — с юмором и эмодзи
- 💼 **Профессиональный** — деловой тон
- 🔥 **Саркастичный** — с иронией
- 💪 **Ободряющий** — позитивный настрой

Выберите лучший вариант и опубликуйте его.

#### 2. Детектор манипуляций 🔍
Напишите текст и нажмите **"🔍 Проверить"**. AI определит наличие манипулятивных техник:
- **Bandwagon** — давление большинства ("все так делают")
- **Authority** — ложная ссылка на авторитеты
- **Fear** — запугивание
- **Urgency** — искусственная срочность
- **Scarcity** — ложный дефицит

При обнаружении вы получите предупреждение и совет, как переформулировать сообщение.

---

## 🛠 Технологический стек

| Компонент | Технологии |
|-----------|------------|
| **Backend API** | Go (Gin), JWT, gorilla/websocket |
| **База данных** | PostgreSQL 15 |
| **Кэш и очереди** | Redis, RabbitMQ |
| **AI сервисы** | Python (FastAPI), DistilBERT, PyTorch |
| **Фронтенд** | HTML5, CSS3, JavaScript (ES6) |
| **Контейнеризация** | Docker Compose / Podman |
| **Хранилище** | MinIO (S3-совместимое) |

---

## 📋 Требования

- **macOS / Linux / Windows** (с WSL2)
- **Docker Desktop** или **Podman Desktop** (рекомендую Podman)
- **Go 1.21+**
- **Python 3.11+**
- **Git**

---

## 🚀 Быстрый старт

### 1. Подготовка окружения (однократно)

```bash
# Установка Go (если не установлен)
brew install go

# Установка Python (если не установлен)
brew install python@3.14

# Проверка
go version
python3 --version
```

### 2. Клонирование и настройка переменных

```bash
git clone https://github.com/StillMortal/yanwit.git
cd yanwit
cp .env.example .env
```

**Отредактируйте** .env **(убедитесь, что параметры верны):**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=yanwit_user
DB_PASSWORD=yanwit_pass
DB_NAME=yanwit

REDIS_ADDR=localhost:6379

RABBITMQ_URL=amqp://yanwit:yanwit_pass@localhost:5672/

API_PORT=8080
JWT_SECRET=yanwit-super-secret-key
```

### 3. Установка зависимостей Go
```bash
cd api
go mod tidy
cd ..
```

### 4. Установка зависимостей Python для AI сервисов
```bash
pip3 install -r ai-services/alternatives/requirements.txt
pip3 install -r ai-services/manipulation/requirements.txt
```

### 5. Запустите инфраструктуру (Docker/Podman)

**Для Podman Desktop (рекомендуется на macOS):**

```bash
podman machine start  # если не запущена
podman-compose up -d
```

**Для Docker Desktop:**

```bash
docker-compose up -d
```

**Что запустится:**
- PostgreSQL (порт 5432)
- Redis (6379)
- RabbitMQ (5672, веб-интерфейс: 15672, логин: yanwit, пароль: yanwit_pass)
- MinIO (9000, консоль: 9001, логин: minioadmin, пароль: minioadmin)

**Проверка:**
```bash
podman ps
# Все контейнеры должны быть в статусе Up
```

### 6. Создание базы данных и миграция (если не создалась автоматически)
```bash
# Создать БД
podman exec -it yanwit-postgres psql -U yanwit_user -d postgres -c "CREATE DATABASE yanwit;" 2>/dev/null || true

# Выполнить миграцию
podman exec -it yanwit-postgres psql -U yanwit_user -d yanwit -f /docker-entrypoint-initdb.d/001_init.sql
```

### 7. Запустите AI сервисы (Python)

**В отдельном терминале — Генератор альтернатив:**
```bash
cd ai-services/alternatives
pip install -r requirements.txt
python3 app.py
```

**Во втором терминале — Детектор манипуляций:**
```bash
cd ai-services/manipulation
pip install -r requirements.txt
python3 app.py
```

**Проверка AI сервисов:**
```bash
curl http://localhost:8002/health
curl http://localhost:8003/health
# Оба должны вернуть {"status":"ok","service":"..."}
```

### 8. Запустите воркера (асинхронная раздача твитов)

**В третьем терминале:**
```bash
cd workers/fanout
go run main.go
```

**Ожидаемый вывод:** Fanout worker started. Waiting for messages...

### 9. Запустите Go API

**В четвёртом терминале:**
```bash
cd api
go run main.go
```

**Ожидаемый вывод:** Yanwit API starting on port 8080

**Проверка API:**
```bash
curl http://localhost:8080/health
# {"database":"connected","service":"yanwit-api","status":"ok"}
```

### 10. Запустите фронтенд

**В пятом терминале:**
```bash
cd frontend
python3 -m http.server 3000
```

**Ожидаемый вывод:** Serving HTTP on :: port 3000

### 11. Откройте приложение в браузере

```text
http://localhost:3000
```
