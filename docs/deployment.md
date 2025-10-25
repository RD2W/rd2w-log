# Deployment Guide

## Локальная разработка

### Предварительные требования
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.25.3+ (опционально для разработки)
- protoc 3.0+ (для генерации protobuf кода)

### Быстрый старт
```bash
# Клонирование репозитория
git clone https://github.com/rd2w/rd2w-log.git
cd rd2w-log

# Генерация protobuf кода
make proto

# Запуск всех сервисов
docker-compose up -d

# Проверка работы
curl http://localhost:8080/api/v1/health
```

### Доступные сервисы
- **API Gateway**: http://localhost:8080
- **PostgreSQL**: http://localhost:5432
- **Redis**: http://localhost:6379
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **Auth Service gRPC**: http://localhost:50051
- **QSO Service gRPC**: http://localhost:50052
- **Analytics Service gRPC**: http://localhost:50053

## Генерация Protobuf кода

### Предварительные требования
```bash
# Установка protoc (Linux/macOS)
# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# macOS
brew install protobuf

# Установка Go плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@vlatest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Генерация кода
```bash
# Генерация всех protobuf файлов
make proto

# Или вручную
./scripts/generate-proto.sh
```

### Структура protobuf файлов
```
proto/
├── auth/
│   └── auth.proto
├── qso/
│   └── qso.proto
├── analytics/
│   └── analytics.proto
├── common/
│   └── common.proto
└── health/
    └── health.proto
```

## Конфигурация

### Environment Variables

#### API Gateway
```bash
API_GATEWAY_HTTP_PORT=8080
API_GATEWAY_GRPC_PORT=50051
AUTH_SERVICE_URL=auth-service:50051
QSO_SERVICE_URL=qso-service:50052
ANALYTICS_SERVICE_URL=analytics-service:50053
JWT_SECRET=your-jwt-secret-key-change-in-production
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

#### Auth Service
```bash
AUTH_DB_HOST=postgres
AUTH_DB_PORT=5432
AUTH_DB_NAME=hamradio
AUTH_DB_USER=hamuser
AUTH_DB_PASSWORD=hampass
JWT_SECRET=your-jwt-secret-key
REDIS_URL=redis://redis:6379
AUTH_SERVICE_PORT=50051
```

#### QSO Service
```bash
QSO_DB_HOST=postgres
QSO_DB_PORT=5432
QSO_DB_NAME=hamradio
QSO_DB_USER=hamuser
QSO_DB_PASSWORD=hampass
REDIS_URL=redis://redis:6379
QSO_SERVICE_PORT=50052
```

#### Analytics Service
```bash
ANALYTICS_DB_HOST=postgres
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_NAME=hamradio
ANALYTICS_DB_USER=hamuser
ANALYTICS_DB_PASSWORD=hampass
ANALYTICS_SERVICE_PORT=50053
```

#### Общие настройки
```bash
LOG_LEVEL=info
ENVIRONMENT=development
APP_NAME=rd2w-log
```

### Docker Compose Configuration

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: hamradio
      POSTGRES_USER: hamuser
      POSTGRES_PASSWORD: hampass
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U hamuser -d hamradio"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5

  auth-service:
    build:
      context: .
      dockerfile: cmd/auth-service/Dockerfile
    environment:
      - AUTH_DB_HOST=postgres
      - AUTH_DB_USER=hamuser
      - AUTH_DB_PASSWORD=hampass
      - JWT_SECRET=your-secret-key-change-in-production
      - AUTH_SERVICE_PORT=50051
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "50051:50051"

  qso-service:
    build:
      context: .
      dockerfile: cmd/qso-service/Dockerfile
    environment:
      - QSO_DB_HOST=postgres
      - QSO_DB_USER=hamuser
      - QSO_DB_PASSWORD=hampass
      - QSO_SERVICE_PORT=50052
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "50052:50052"

  analytics-service:
    build:
      context: .
      dockerfile: cmd/analytics-service/Dockerfile
    environment:
      - ANALYTICS_DB_HOST=postgres
      - ANALYTICS_DB_USER=hamuser
      - ANALYTICS_DB_PASSWORD=hampass
      - ANALYTICS_SERVICE_PORT=50053
    depends_on:
      postgres:
        condition: service_healthy
    ports:
      - "50053:50053"

  api-gateway:
    build:
      context: .
      dockerfile: cmd/api-gateway/Dockerfile
    environment:
      - AUTH_SERVICE_URL=auth-service:50051
      - QSO_SERVICE_URL=qso-service:50052
      - ANALYTICS_SERVICE_URL=analytics-service:50053
      - API_GATEWAY_HTTP_PORT=8080
    depends_on:
      - auth-service
      - qso-service
      - analytics-service
    ports:
      - "8080:8080"

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./deployments/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./deployments/grafana/dashboards:/etc/grafana/provisioning/dashboards

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
```

## Production Deployment

### Требования для production
- Kubernetes cluster или cloud platform
- External PostgreSQL database
- SSL certificates
- Load balancer
- Redis cluster

### Kubernetes Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: your-registry/ham-radio-auth:latest
        ports:
        - containerPort: 50051
        env:
        - name: AUTH_DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secret
              key: jwt-secret
        livenessProbe:
          exec:
            command: ["grpc_health_probe", "-addr=:50051"]
          initialDelaySeconds: 10
          periodSeconds: 5
        readinessProbe:
          exec:
            command: ["grpc_health_probe", "-addr=:50051"]
          initialDelaySeconds: 5
          periodSeconds: 2
```

### Health Checks

Все сервисы предоставляют health endpoints:

#### gRPC Health Checks
```bash
# Проверка здоровья через grpcurl
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Проверка конкретного сервиса
grpcurl -plaintext -d '{"service": "auth.AuthService"}' localhost:50051 grpc.health.v1.Health/Check

# Использование grpc_health_probe
grpc_health_probe -addr=localhost:50051
```

#### HTTP Health Checks
```bash
# Основной health check
curl http://localhost:8080/health

# Детальная проверка
curl http://localhost:8080/health/detailed

# Проверка готовности
curl http://localhost:8080/health/ready

# Проверка живучести
curl http://localhost:8080/health/live
```

### Мониторинг

#### Prometheus Configuration
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'auth-service'
    static_configs:
      - targets: ['auth-service:50051']
    metrics_path: /metrics

  - job_name: 'qso-service'
    static_configs:
      - targets: ['qso-service:50052']
    metrics_path: /metrics

  - job_name: 'analytics-service'
    static_configs:
      - targets: ['analytics-service:50053']
    metrics_path: /metrics

  - job_name: 'api-gateway'
    static_configs:
      - targets: ['api-gateway:8080']
    metrics_path: /metrics
```

#### Grafana Dashboards
Импортируйте готовые дашборды из `deployments/grafana/dashboards/`

## Миграции базы данных

```bash
# Запуск миграций
docker-compose run --rm migrate -path=/migrations -database "$DB_URL" up

# Откат миграций
docker-compose run --rm migrate -path=/migrations -database "$DB_URL" down

# Создание новой миграции
docker-compose run --rm migrate create -ext sql -dir /migrations -seq migration_name
```

## Backup и Recovery

### PostgreSQL Backup
```bash
# Создание бэкапа
docker-compose exec postgres pg_dump -U hamuser hamradio > backup_$(date +%Y%m%d_%H%M%S).sql

# Восстановление из бэкапа
docker-compose exec -T postgres psql -U hamuser hamradio < backup.sql

# Автоматические бэкапы с retention
docker-compose exec postgres pg_dump -U hamuser -Fc hamradio > backup_$(date +%Y%m%d_%H%M%S).dump
```

### Redis Backup
```bash
# Создание RDB файла
docker-compose exec redis redis-cli SAVE

# Копирование файла
docker cp rd2w-log_redis_1:/data/dump.rdb dump_$(date +%Y%m%d_%H%M%S).rdb

# Настройка автоматических бэкапов в redis.conf
save 900 1
save 300 10
save 60 10000
```

## Troubleshooting

### Проверка логов
```bash
# Все сервисы
docker-compose logs

# Конкретный сервис
docker-compose logs auth-service

# Логи в реальном времени
docker-compose logs -f qso-service

# Логи с фильтрацией по уровню
docker-compose logs --tail=100 auth-service | grep "ERROR"
```

### Проверка здоровья
```bash
# Проверка всех сервисов
docker-compose ps

# Проверка базы данных
docker-compose exec postgres pg_isready -U hamuser

# Проверка Redis
docker-compose exec redis redis-cli ping

# Проверка gRPC сервисов
grpcurl -plaintext localhost:50051 list
```

### Распространенные проблемы

1. **Port conflicts** - проверить что порты 8080, 5432, 6379, 50051-50053 свободны
2. **Memory issues** - увеличьте лимиты памяти в Docker
3. **Database connection** - проверить credentials в environment variables
4. **Protobuf generation** - убедиться что protoc и плагины установлены правильно
5. **JWT secret** - использовать надежный JWT секрет в production
