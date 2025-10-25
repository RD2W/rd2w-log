# Техническое задание: Ham Radio QSO Journal

## Оглавление
1. [Введение](#введение)
2. [Технический стек](#технический-стек)
3. [Бизнес-сущности](#бизнес-сущности)
4. [Архитектура системы](#архитектура-системы)
5. [Структура проекта](#структура-проекта)
6. [Protobuf спецификации](#protobuf-спецификации)
7. [Базы данных](#базы-данных)
8. [Безопасность](#безопасность)
9. [Тестирование](#тестирование)
10. [Мониторинг](#мониторинг)
11. [План разработки](#план-разработки)
12. [Критерии успеха](#критерии-успеха)

## Введение

### Описание проекта
Ham Radio QSO Journal - это микросервисная система для ведения журнала радиосвязей (QSO) с расширенной аналитикой и поддержкой мобильных клиентов.

### Цели проекта
- Создать современную, масштабируемую систему для управления QSO
- Обеспечить удобный интерфейс для радиолюбителей
- Предоставить аналитику и отчетность
- Поддержать стандартные форматы (ADIF)

### Целевая аудитория
- Радиолюбители, участвующие в соревнованиях
- DX-экспедиции
- Любители повседневной радиосвязи

## Технический стек

### Backend
- **Язык**: Go 1.25.3+
- **Архитектура**: Clean Architecture
- **Коммуникация**: gRPC + Protocol Buffers
- **API Gateway**: gRPC Gateway (REST/JSON)
- **HTTP Framework**: Gin

### Базы данных
- **Основная БД**: PostgreSQL 15+
- **Кэш**: Redis 7+
- **Миграции**: Goose

### Инфраструктура
- **Контейнеризация**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Мониторинг**: Prometheus + Grafana
- **Логирование**: Structured logging (Zap)
- **Health Checks**: gRPC Health Protocol

## Бизнес-сущности

### Пользователь (User)
```protobuf
message User {
  string id = 1;
  string callsign = 2;
  string email = 3;
  string name = 4;
  string location = 5;
  string timezone = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp updated_at = 8;
}
```

### Радиосвязь (QSO)
```protobuf
message QSO {
  string id = 1;
  string user_id = 2;
  string callsign = 3;
  google.protobuf.Timestamp time_on = 4;
  google.protobuf.Timestamp time_off = 5;
  string frequency = 6;
  string band = 7;
  string mode = 8;
  string rst_sent = 9;
  string rst_received = 10;
  string grid_square = 11;
  string operator = 12;
  string station_callsign = 13;
  string contest = 14;
  string notes = 15;
  QSLInfo qsl = 16;
  google.protobuf.Timestamp created_at = 17;
  google.protobuf.Timestamp updated_at = 18;
}
```

### QSL информация (QSLInfo)
```protobuf
message QSLInfo {
  bool sent = 1;
  bool received = 2;
  string sent_via = 3;
  string received_via = 4;
  google.protobuf.Timestamp sent_date = 5;
  google.protobuf.Timestamp received_date = 6;
}
```

### Станция (Station)
```protobuf
message Station {
  string id = 1;
  string user_id = 2;
  string callsign = 3;
  string name = 4;
  string location = 5;
  string grid_square = 6;
  repeated string equipment = 7;
  repeated string antennas = 8;
  repeated string supported_modes = 9;
  repeated string supported_bands = 10;
}
```

### Статистика (Statistics)
```protobuf
message Statistics {
  int32 total_qso = 1;
  map<string, int32> qso_by_band = 2;
  map<string, int32> qso_by_mode = 3;
  map<string, int32> qso_by_continent = 4;
  repeated int32 last_30_days = 5;
  int32 unique_countries = 6;
  int32 unique_grids = 7;
}
```

## Архитектура системы

### Микросервисная архитектура

#### 1. Сервис аутентификации (Auth Service)
**Ответственность**: Управление пользователями, аутентификация, авторизация

**gRPC методы**:
- `Register` - регистрация нового пользователя
- `Login` - вход в систему
- `ValidateToken` - валидация JWT токена
- `RefreshToken` - обновление токена
- `UpdateProfile` - обновление профиля пользователя
- `GetProfile` - получение профиля пользователя

#### 2. Сервис QSO (QSO Service)
**Ответственность**: Управление журналом радиосвязей

**gRPC методы**:
- `CreateQSO`, `GetQSO`, `UpdateQSO`, `DeleteQSO` - CRUD операции
- `ListQSOs` - потоковый список с пагинацией и фильтрацией
- `SearchQSOs` - поиск по критериям
- `ImportADIF`, `ExportADIF` - работа с ADIF форматом

#### 3. Сервис аналитики (Analytics Service)
**Ответственность**: Статистика и аналитика

**gRPC методы**:
- `GetStatistics` - общая статистика по связям
- `GetExtendedStatistics` - расширенная статистика с метриками роста
- `GetActivityHeatmap` - тепловая карта активности
- `GenerateReports` - генерация отчетов в различных форматах

#### 4. API Gateway
**Ответственность**: Единая точка входа для REST API

**Функциональность**:
- Маршрутизация HTTP запросов к gRPC сервисам
- Аутентификация и авторизация
- Преобразование JSON ↔ Protobuf
- Rate limiting
- CORS

## Структура проекта

### Общая структура проекта:
```
cmd/                    # Точки входа приложений
├── auth-service/       # Сервис аутентификации
├── qso-service/        # Сервис QSO
├── analytics-service/  # Сервис аналитики
└── api-gateway/        # API Gateway

internal/               # Внутренние пакеты
├── auth/               # Домен аутентификации
│   ├── domain/         # Бизнес-сущности
│   │   ├── user.go
│   │   ├── session.go
│   │   ├── errors.go
│   │   └── types.go
│   ├── usecase/        # Бизнес-логика
│   ├── delivery/       # gRPC handlers
│   └── repository/     # Репозитории (PostgreSQL, Redis)
├── qso/                # Домен QSO
│   ├── domain/
│   │   ├── qso.go
│   │   ├── qsl_info.go
│   │   ├── station.go
│   │   ├── filter.go
│   │   ├── errors.go
│   │   └── types.go
│   ├── usecase/
│   ├── delivery/
│   └── repository/
├── analytics/          # Домен аналитики
│   ├── domain/
│   │   ├── statistics.go
│   │   ├── activity.go
│   │   ├── report.go
│   │   ├── errors.go
│   │   └── types.go
│   ├── usecase/
│   ├── delivery/
│   └── repository/
└── shared/             # Общие компоненты
    ├── config/         # Конфигурация
    ├── database/       # Подключение к БД
    ├── jwt/            # JWT утилиты
    ├── logger/         # Логирование
    └── grpc_client/    # gRPC клиенты

pkg/                    # Переиспользуемые пакеты
├── proto/              # Сгенерированные protobuf файлы
├── adif/               # Парсер ADIF формата
├── validator/          # Валидация данных
├── cache/              # Кэширование
└── metrics/            # Метрики Prometheus

proto/                  # Protobuf определения
├── auth/               # Auth service protos
├── qso/                # QSO service protos
├── analytics/          # Analytics service protos
├── common/             # Общие сообщения
└── health/             # Health check protos

deployments/            # Деплой и инфраструктура
├── docker-compose.yml  # Локальный запуск
├── prometheus/         # Конфигурация Prometheus
└── grafana/            # Дашборды Grafana

configs/                # Файлы конфигурации
├── auth-config.yaml
├── qso-config.yaml
└── analytics-config.yaml

scripts/                # Вспомогательные скрипты
└── generate-proto.sh   # Генерация protobuf кода
```

## Protobuf спецификации

### Сервис аутентификации (Auth Service)
- `Register` - регистрация нового пользователя
- `Login` - аутентификация пользователя
- `ValidateToken` - валидация JWT токена
- `RefreshToken` - обновление access токена
- `UpdateProfile` - обновление профиля пользователя
- `GetProfile` - получение профиля пользователя

### Сервис QSO (QSO Service)
- `CreateQSO`, `GetQSO`, `UpdateQSO`, `DeleteQSO` - CRUD операции
- `ListQSOs` - потоковый список с пагинацией и фильтрацией
- `SearchQSOs` - поиск по позывным и заметкам
- `ImportADIF`, `ExportADIF` - работа с ADIF форматом

### Сервис аналитики (Analytics Service)
- `GetStatistics` - базовая статистика
- `GetExtendedStatistics` - расширенная статистика с метриками роста
- `GetActivityHeatmap` - тепловая карта активности
- `GenerateReports` - генерация отчетов в различных форматах

### Общие сообщения
- `Error` - стандартизированная ошибка
- `Pagination` - информация о пагинации
- `TimeRange` - временной диапазон

### Protobuf структура
```
proto/                    # Protobuf определения
├── auth/                 # Auth service protos
│   └── auth.proto
├── qso/                  # QSO service protos
│   └── qso.proto
├── analytics/            # Analytics service protos
│   └── analytics.proto
├── common/               # Общие сообщения
│   └── common.proto
└── health/               # Health check protos
    └── health.proto

pkg/proto/                # Сгенерированный Go код
├── auth/                 # Сгенерированные файлы auth
├── qso/                  # Сгенерированные файлы qso
├── analytics/            # Сгенерированные файлы analytics
├── common/               # Сгенерированные общие файлы
└── health/               # Сгенерированные health файлы
```

### Скрипты генерации
- `scripts/generate-proto.sh` - скрипт генерации Go кода
- `Makefile` - цели для генерации protobuf

## Базы данных

### PostgreSQL схемы

#### Пользователи:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    callsign VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    location VARCHAR(100),
    timezone VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### QSO записи:
```sql
CREATE TABLE qso_records (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    callsign VARCHAR(20) NOT NULL,
    time_on TIMESTAMPTZ NOT NULL,
    time_off TIMESTAMPTZ,
    frequency DECIMAL(10,6),
    band VARCHAR(10),
    mode VARCHAR(10),
    rst_sent VARCHAR(10),
    rst_received VARCHAR(10),
    grid_square VARCHAR(10),
    operator VARCHAR(50),
    station_callsign VARCHAR(20),
    contest VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### QSL информация:
```sql
CREATE TABLE qsl_info (
    id UUID PRIMARY KEY,
    qso_id UUID REFERENCES qso_records(id),
    sent BOOLEAN DEFAULT FALSE,
    received BOOLEAN DEFAULT FALSE,
    sent_via VARCHAR(20),
    received_via VARCHAR(20),
    sent_date TIMESTAMPTZ,
    received_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Станции:
```sql
CREATE TABLE stations (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    callsign VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100),
    grid_square VARCHAR(10),
    equipment TEXT[], -- Array of equipment
    antennas TEXT[],  -- Array of antennas
    supported_modes TEXT[], -- Array of supported modes
    supported_bands TEXT[], -- Array of supported bands
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Индексы:
```sql
CREATE INDEX idx_qso_user_id ON qso_records(user_id);
CREATE INDEX idx_qso_time_on ON qso_records(time_on);
CREATE INDEX idx_qso_callsign ON qso_records(callsign);
CREATE INDEX idx_qso_band ON qso_records(band);
CREATE INDEX idx_qso_mode ON qso_records(mode);
CREATE INDEX idx_qso_grid_square ON qso_records(grid_square);
CREATE INDEX idx_stations_user_id ON stations(user_id);
```

## Безопасность

### Аутентификация
- JWT tokens с expiration (15 минут для access, 7 дней для refresh)
- Refresh token механизм с rotation
- HTTPS для всех коммуникаций
- Password hashing с bcrypt

### Защита данных
- Валидация входных данных
- SQL injection protection (prepared statements)
- Rate limiting (по IP и пользователю)
- CORS настройки
- Input sanitization

### Безопасность API
- API ключи для внешних интеграций
- Request signing
- Timestamp validation для предотвращения replay attacks

## Тестирование

### Стратегия тестирования
- **Юнит-тесты**: Бизнес-логика, доменные сущности (покрытие > 80%)
- **Интеграционные тесты**: Репозитории, gRPC handlers
- **E2E тесты**: Полные сценарии использования API
- **Нагрузочное тестирование**: Проверка производительности

### Метрики качества
- Покрытие кода > 80%
- Отсутствие гонок данных
- Статический анализ (golangci-lint)
- Security scanning (gosec)

## Мониторинг

### Prometheus метрики
- Количество QSO по часам/дням
- Время ответа gRPC методов
- Количество ошибок по типам
- Использование кэша Redis
- Database connection pool metrics
- gRPC stream metrics

### Grafana дашборды
- Общий обзор системы
- QSO активность и статистика
- Производительность сервисов
- Бизнес-метрики
- Error rates и latency

### Health Checks
- gRPC health checking protocol
- HTTP health endpoints
- Readiness и liveness probes
- Dependency health checks (DB, Redis)

## План разработки

### Фаза 1: Базовая функциональность
- [x] Настройка проекта и CI/CD
- [x] Проектирование архитектуры и protobuf
- [ ] Реализация Auth Service
- [ ] Базовые CRUD операции для QSO
- [ ] Настройка gRPC Gateway

### Фаза 2: Расширенная функциональность
- [ ] Реализация Analytics Service
- [ ] ADIF импорт/экспорт
- [ ] Расширенная фильтрация QSO
- [ ] Кэширование с Redis
- [ ] Управление станциями

### Фаза 3: Продвинутые возможности
- [ ] Расширенная аналитика и отчеты
- [ ] Real-time уведомления
- [ ] Интеграция с внешними API
- [ ] Оптимизация производительности

### Фаза 4: Production готовность
- [ ] Мониторинг и метрики
- [ ] Нагрузочное тестирование
- [ ] Документация API
- [ ] Security audit
- [ ] Production деплой

## Критерии успеха

### Технические критерии
- [ ] Все микросервисы запускаются через Docker Compose
- [ ] CI/CD pipeline успешно проходит все этапы
- [ ] Покрытие кода тестами > 80%
- [ ] Отсутствие critical security issues
- [ ] Генерация protobuf кода работает корректно

### Функциональные критерии
- [ ] Полный цикл создания и управления QSO
- [ ] Работающая аутентификация и авторизация
- [ ] Импорт/экспорт ADIF файлов
- [ ] Базовая и расширенная аналитика
- [ ] Управление станциями

### Бизнес-критерии
- [ ] Система решает реальные проблемы радиолюбителей
- [ ] Интерфейс понятен целевой аудитории
- [ ] Производительность удовлетворяет требованиям
- [ ] Система масштабируется под нагрузку
- 