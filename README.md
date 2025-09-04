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
env: "localv"
addr: "0.0.0.0:44032"
services:
  sso:
    endpoint: "sso-app:44043"
    timeout: "5s"
  task:
    endpoint: "task-app:44045"
    timeout: "5s"
```

3. **Запустите сервис**
```bash
# Локальная разработка
go run cmd/api-gateway/main.go -config=config/local.yaml

# Или через Docker
docker compose up --build
```

## 📡 API Endpoints

### Аутентификация
```http
POST /auth/login
POST /auth/register  
POST /auth/refresh
```

### Управление задачами
```http
POST   /tasks               # Создать новую задачу
GET    /tasks/{id}          # Получить задачу по ID
PUT    /tasks/{id}          # Обновить задачу
DELETE /tasks/{id}          # Удалить задачу
PATCH  /tasks/{id}/status   # Изменить статус задачи
```

## 🛡️ Middleware

### Authentication Middleware
Проверяет JWT токены через SSO сервис для защищенных endpoints.

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
