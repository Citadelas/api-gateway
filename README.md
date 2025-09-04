# API Gateway - Citadelas

**API Gateway** для микросервисной архитектуры Citadelas. Обеспечивает единую точку входа для всех клиентских запросов и маршрутизирует их к соответствующим микросервисам.

## 🏗️ Архитектура

API Gateway служит центральным узлом для:
- **Аутентификация и авторизация** через SSO сервис
- **Управление задачами** через Task сервис  
- **Маршрутизация запросов** к backend сервисам
- **Логирование и мониторинг** всех запросов

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.24+
- Docker & Docker Compose
- Запущенные SSO и Task сервисы

### Установка и запуск

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/Citadelas/api-gateway.git
cd api-gateway
```

2. **Создайте конфигурационный файл**
Создайте файл `config/local.yaml`:
```yaml
env: "local"
server:
  port: 8080
  timeout: 30s

services:
  sso:
    url: "localhost:44043"
    timeout: 5s
  task:
    url: "localhost:44044" 
    timeout: 5s

logging:
  level: "debug"
  format: "json"
```

3. **Запустите сервис**
```bash
# Локальная разработка
go run cmd/api-gateway/main.go -config=config/local.yaml

# Или через Docker
docker build -t citadelas/api-gateway .
docker run -p 8080:8080 -v $(pwd)/config:/app/config citadelas/api-gateway
```

## 📡 API Endpoints

### Аутентификация
```http
POST /auth/login
POST /auth/register  
POST /auth/refresh
GET  /auth/me
```

### Управление задачами
```http
GET    /tasks               # Получить все задачи пользователя
POST   /tasks               # Создать новую задачу
GET    /tasks/{id}          # Получить задачу по ID
PUT    /tasks/{id}          # Обновить задачу
DELETE /tasks/{id}          # Удалить задачу
PATCH  /tasks/{id}/status   # Изменить статус задачи
```

### Системные endpoints
```http
GET /health    # Health check
GET /ready     # Readiness probe
GET /metrics   # Prometheus metrics
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|-----------|----------|--------------|
| `CONFIG_PATH` | Путь к файлу конфигурации | `./config/local.yaml` |
| `PORT` | Порт сервера | `8080` |
| `SSO_SERVICE_URL` | URL SSO сервиса | `localhost:44043` |
| `TASK_SERVICE_URL` | URL Task сервиса | `localhost:44044` |

### Структура конфигурации

```yaml
env: "local|dev|prod"
server:
  port: 8080
  timeout: "30s"
  read_timeout: "10s"
  write_timeout: "10s"

services:
  sso:
    url: "sso-service:44043"
    timeout: "5s"
    max_retries: 3
  task:
    url: "task-service:44044"
    timeout: "5s" 
    max_retries: 3

cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "PATCH"]
  allowed_headers: ["*"]

logging:
  level: "info"
  format: "json"
```

## 🛡️ Middleware

### Authentication Middleware
Проверяет JWT токены через SSO сервис для защищенных endpoints.

### CORS Middleware  
Настраивает Cross-Origin Resource Sharing для фронтенд приложений.

### Logging Middleware
Логирует все входящие запросы с метриками производительности.

### Rate Limiting Middleware *(Планируется)*
Защита от DDoS атак и злоупотреблений API.

## 🏃‍♂️ Разработка

### Структура проекта

```
api-gateway/
├── cmd/
│   └── api-gateway/
│       └── main.go          # Точка входа
├── internal/
│   ├── app/                 # Инициализация приложения
│   ├── config/              # Управление конфигурацией
│   ├── handlers/            # HTTP обработчики
│   │   ├── auth/           # Аутентификация endpoints
│   │   └── tasks/          # Задачи endpoints
│   ├── middleware/          # HTTP middleware
│   ├── services/            # Бизнес-логика и gRPC клиенты
│   └── lib/
│       └── logger/         # Логирование
├── config/                 # Конфигурационные файлы
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

### Добавление нового endpoint

1. Создайте handler в соответствующем пакете
2. Добавьте маршрут в `internal/app/app.go`
3. При необходимости добавьте middleware
4. Обновите документацию API

## 🐳 Docker

### Dockerfile
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-gateway cmd/api-gateway/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config ./config
CMD ["./api-gateway"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  api-gateway:
    build: .
    ports:
      - "8080:8080"
    environment:
      - CONFIG_PATH=/app/config/local.yaml
      - SSO_SERVICE_URL=sso:44043
      - TASK_SERVICE_URL=task:44044
    depends_on:
      - sso
      - task
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## 📊 Мониторинг и логирование

### Health Checks
- `/health` - общее состояние сервиса
- `/ready` - готовность к обслуживанию запросов

### Metrics *(Планируется)*
- HTTP request duration
- Request count по endpoints
- Error rate по статус кодам
- gRPC клиентские метрики

### Logging
Структурированные логи в JSON формате включают:
- Request ID для трейсинга
- HTTP метод и путь
- Время выполнения запроса
- Статус код ответа
- User ID (если аутентифицирован)

## 🚧 Roadmap

### v1.1
- [ ] Rate limiting middleware
- [ ] Circuit breaker для gRPC клиентов
- [ ] Prometheus metrics
- [ ] Request validation middleware

### v1.2  
- [ ] API versioning
- [ ] Caching layer с Redis
- [ ] Load balancing для backend сервисов
- [ ] Distributed tracing

### v2.0
- [ ] GraphQL gateway
- [ ] Websocket поддержка
- [ ] API documentation автогенерация
- [ ] Admin panel для мониторинга

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.

## 🔗 Related Services

- [SSO Service](https://github.com/Citadelas/sso) - Аутентификация и авторизация
- [Task Service](https://github.com/Citadelas/task) - Управление задачами
- [Protos](https://github.com/Citadelas/protos) - Protocol Buffers контракты

---

**Maintainer**: [muerewa](https://github.com/muerewa)  
**Organization**: [Citadelas](https://github.com/Citadelas)
