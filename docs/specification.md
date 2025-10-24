# Техническое задание: Ham Radio QSO Journal

## Оглавление
1. [Введение](#введение)
2. [Технический стек](#технический-стек)
3. [Бизнес-сущности](#бизнес-сущности)
4. [Архитектура системы](#архитектура-системы)
5. [Структура проекта](#структура-проекта)
6. [Базы данных](#базы-данных)
7. [Безопасность](#безопасность)
8. [Тестирование](#тестирование)
9. [Мониторинг](#мониторинг)
10. [План разработки](#план-разработки)
11. [Критерии успеха](#критерии-успеха)

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

### Базы данных
- **Основная БД**: PostgreSQL 15+
- **Кэш**: Redis 7+
- **Миграции**: Goose

### Инфраструктура
- **Контейнеризация**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Мониторинг**: Prometheus + Grafana
- **Логирование**: Structured logging

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

#### 2. Сервис QSO (QSO Service)
**Ответственность**: Управление журналом радиосвязей

**gRPC методы**:
- `CreateQSO`, `GetQSO`, `UpdateQSO`, `DeleteQSO`
- `ListQSOs` - потоковый список с фильтрацией
- `SearchQSOs` - поиск по критериям
- `ImportADIF`, `ExportADIF` - работа с ADIF форматом

#### 3. Сервис аналитики (Analytics Service)
**Ответственность**: Базовая статистика и аналитика

**gRPC методы**:
- `GetStatistics` - общая статистика по связям
- `GetActivityHeatmap` - тепловая карта активности
- `GenerateReports` - генерация отчетов

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
│   ├── domain/         # Бизнес-сущности и интерфейсы
│   ├── usecase/        # Бизнес-логика
│   ├── delivery/       # gRPC handlers
│   └── repository/     # Репозитории (PostgreSQL, Redis)
├── qso/                # Домен QSO
├── analytics/          # Домен аналитики
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
└── analytics/          # Analytics service protos

deployments/            # Деплой и инфраструктура
├── docker-compose.yml  # Локальный запуск
├── prometheus/         # Конфигурация Prometheus
└── grafana/            # Дашборды Grafana

configs/                # Файлы конфигурации
├── auth-config.yaml
├── qso-config.yaml
└── analytics-config.yaml
```

### Детальное описание структуры:

#### cmd/ - Точки входа
- `auth-service/` - основной файл сервиса аутентификации
- `qso-service/` - основной файл сервиса QSO
- `analytics-service/` - основной файл сервиса аналитики
- `api-gateway/` - основной файл API Gateway

#### internal/ - Внутренняя архитектура
Каждый сервис имеет следующую структуру:

**Домен аутентификации (auth/)**:
- `domain/` - бизнес-сущности (User, Session) и интерфейсы репозиториев
- `usecase/` - бизнес-логика (регистрация, аутентификация, валидация)
- `delivery/` - gRPC handlers и interceptors
- `repository/` - реализации репозиториев (PostgreSQL, Redis)

**Домен QSO (qso/)**:
- `domain/` - сущности QSO, Station, QSLInfo
- `usecase/` - логика управления QSO, поиска, импорта/экспорта
- `delivery/` - gRPC handlers для QSO операций
- `repository/` - репозитории для работы с QSO данными

**Домен аналитики (analytics/)**:
- `domain/` - сущности отчетов и статистики
- `usecase/` - логика агрегации данных и генерации отчетов
- `delivery/` - gRPC handlers для аналитических запросов
- `repository/` - репозитории для аналитических данных

**Общие компоненты (shared/)**:
- `config/` - загрузка и валидация конфигурации
- `database/` - подключения к PostgreSQL и Redis
- `jwt/` - генерация и валидация JWT токенов
- `logger/` - структурированное логирование
- `grpc_client/` - клиенты для межсервисного взаимодействия

#### pkg/ - Переиспользуемые пакеты
- `proto/` - сгенерированные Go файлы из .proto определений
- `adif/` - парсер и генератор ADIF формата
- `validator/` - валидация входных данных
- `cache/` - абстракции для кэширования
- `metrics/` - метрики Prometheus

#### Дополнительные директории
- `proto/` - исходные .proto файлы для gRPC сервисов
- `deployments/` - конфигурации для развертывания и мониторинга
- `configs/` - YAML конфигурационные файлы для сервисов

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

#### Индексы:
```sql
CREATE INDEX idx_qso_user_id ON qso_records(user_id);
CREATE INDEX idx_qso_time_on ON qso_records(time_on);
CREATE INDEX idx_qso_callsign ON qso_records(callsign);
CREATE INDEX idx_qso_band ON qso_records(band);
CREATE INDEX idx_qso_mode ON qso_records(mode);
```

## Безопасность

### Аутентификация
- JWT tokens с expiration
- Refresh token механизм
- HTTPS для всех коммуникаций

### Защита данных
- Хеширование паролей (bcrypt)
- Валидация входных данных
- SQL injection protection
- Rate limiting
- CORS настройки

## Тестирование

### Стратегия тестирования
- **Юнит-тесты**: Бизнес-логика, доменные сущности
- **Интеграционные тесты**: Репозитории, gRPC handlers
- **E2E тесты**: Полные сценарии использования

### Метрики качества
- Покрытие кода > 80%
- Отсутствие гонок данных
- Статический анализ (golangci-lint)

## Мониторинг

### Prometheus метрики
- Количество QSO по часам/дням
- Время ответа gRPC методов
- Количество ошибок по типам
- Использование кэша Redis

### Grafana дашборды
- Общий обзор системы
- QSO активность
- Производительность сервисов
- Бизнес-метрики

## План разработки

### Фаза 1: Базовая функциональность
- [ ] Настройка проекта и CI/CD
- [ ] Реализация Auth Service
- [ ] Базовые CRUD операции для QSO
- [ ] Настройка gRPC Gateway

### Фаза 2: Расширенная функциональность
- [ ] Реализация Analytics Service
- [ ] ADIF импорт/экспорт
- [ ] Расширенная фильтрация QSO
- [ ] Кэширование с Redis

### Фаза 3: Продвинутые возможности
- [ ] Управление станциями
- [ ] Расширенная аналитика
- [ ] Real-time уведомления
- [ ] Оптимизация производительности

### Фаза 4: Production готовность
- [ ] Мониторинг и метрики
- [ ] Нагрузочное тестирование
- [ ] Документация API
- [ ] Production деплой

## Критерии успеха

### Технические критерии
- [ ] Все микросервисы запускаются через Docker Compose
- [ ] CI/CD pipeline успешно проходит все этапы
- [ ] Покрытие кода тестами > 80%
- [ ] Отсутствие critical security issues

### Функциональные критерии
- [ ] Полный цикл создания и управления QSO
- [ ] Работающая аутентификация и авторизация
- [ ] Импорт/экспорт ADIF файлов
- [ ] Базовая аналитика и статистика

### Бизнес-критерии
- [ ] Система решает реальные проблемы радиолюбителей
- [ ] Интерфейс понятен целевой аудитории
- [ ] Производительность удовлетворяет требованиям
