# Deployment Guide

## Локальная разработка

### Предварительные требования
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.25.3+ (опционально для разработки)

### Быстрый старт
```bash
# Клонирование репозитория
git clone https://github.com/rd2w/rd2w-log.git
cd ham-radio-qso-journal

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

## Конфигурация

### Environment Variables

#### Auth Service
```bash
AUTH_DB_HOST=postgres
AUTH_DB_PORT=5432
AUTH_DB_NAME=hamradio
AUTH_DB_USER=hamuser
AUTH_DB_PASSWORD=hampass
JWT_SECRET=your-jwt-secret-key
REDIS_URL=redis://redis:6379
```

#### QSO Service
```bash
QSO_DB_HOST=postgres
QSO_DB_PORT=5432
QSO_DB_NAME=hamradio
QSO_DB_USER=hamuser
QSO_DB_PASSWORD=hampass
REDIS_URL=redis://redis:6379
```

#### Analytics Service
```bash
ANALYTICS_DB_HOST=postgres
ANALYTICS_DB_PORT=5432
ANALYTICS_DB_NAME=hamradio
ANALYTICS_DB_USER=hamuser
ANALYTICS_DB_PASSWORD=hampass
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
```

### Health Checks

Все сервисы предоставляют health endpoints:
- gRPC: `grpc.health.v1.Health/Check`
- HTTP: `/health`

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
```

#### Grafana Dashboards
Импортируйте готовые дашборды из `deployments/grafana/dashboards/`

## Миграции базы данных

```bash
# Запуск миграций
docker-compose run --rm migrate -path=/migrations -database "$DB_URL" up

# Откат миграций
docker-compose run --rm migrate -path=/migrations -database "$DB_URL" down
```

## Backup и Recovery

### PostgreSQL Backup
```bash
# Создание бэкапа
docker-compose exec postgres pg_dump -U hamuser hamradio > backup.sql

# Восстановление из бэкапа
docker-compose exec -T postgres psql -U hamuser hamradio < backup.sql
```

### Redis Backup
```bash
# Создание RDB файла
docker-compose exec redis redis-cli SAVE

# Копирование файла
docker cp ham-radio-qso-journal_redis_1:/data/dump.rdb .
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
```

### Проверка здоровья
```bash
# Проверка всех сервисов
docker-compose ps

# Проверка базы данных
docker-compose exec postgres pg_isready -U hamuser

# Проверка Redis
docker-compose exec redis redis-cli ping
```

### Распространенные проблемы

1. **Port conflicts** - проверить что порты 8080, 5432, 6379 свободны
2. **Memory issues** - увеличьте лимиты памяти в Docker
3. **Database connection** - проверить credentials в environment variables
