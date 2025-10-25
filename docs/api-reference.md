# API Reference

## Версия API
Текущая версия: v1.0.0

## Общая информация

Базовый URL: `http://localhost:8080/api/v1`

### Аутентификация
Большинство endpoints требуют JWT аутентификации, поэтому необходимо добавлять токен в заголовок:
```
Authorization: Bearer <your_jwt_token>
```

### Формат дат
Все даты передаются в формате ISO 8601: `YYYY-MM-DDTHH:MM:SSZ`

### Пагинация
Все списковые endpoints поддерживают пагинацию через query parameters:
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице (максимум 1000)

### Ошибки
Все ошибки возвращаются в формате:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Описание ошибки",
    "details": {
      "field": "дополнительная информация"
    }
  }
}
```

### Статусы операций
- `200` - Успешный запрос
- `201` - Успешное создание
- `400` - Неверные параметры запроса
- `401` - Требуется аутентификация
- `403` - Недостаточно прав
- `404` - Ресурс не найден
- `409` - Конфликт данных (например, пользователь уже существует)
- `500` - Внутренняя ошибка сервера

## Health Check Endpoints

### Проверка здоровья сервиса
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T14:30:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "external_api": "healthy"
  }
}
```

## Auth Endpoints

### Регистрация пользователя
```http
POST /auth/register
Content-Type: application/json

{
  "callsign": "RD2W",
  "email": "user@example.com",
  "password": "securepassword",
  "name": "Maxim Kru",
  "location": "Kursk, Russia",
  "timezone": "Europe/Moscow"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "callsign": "RD2W",
    "email": "user@example.com",
    "name": "Maxim Kru",
    "location": "Kursk, Russia",
    "timezone": "Europe/Moscow",
    "created_at": "2024-01-15T14:30:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  },
  "access_token": "jwt_token",
  "refresh_token": "refresh_token"
}
```

### Вход в систему
```http
POST /auth/login
Content-Type: application/json

{
  "callsign": "RD2W",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "callsign": "RD2W",
    "email": "user@example.com",
    "name": "Maxim Kru",
    "location": "Kursk, Russia",
    "timezone": "Europe/Moscow"
  },
  "access_token": "jwt_token",
  "refresh_token": "refresh_token"
}
```

### Обновление токена
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token"
}
```

**Response:**
```json
{
  "access_token": "new_jwt_token",
  "refresh_token": "new_refresh_token"
}
```

### Получение профиля пользователя
```http
GET /user/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "callsign": "RD2W",
    "email": "user@example.com",
    "name": "Maxim Kru",
    "location": "Kursk, Russia",
    "timezone": "Europe/Moscow",
    "created_at": "2024-01-15T14:30:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  }
}
```

### Обновление профиля пользователя
```http
PUT /user/profile
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "New Name",
  "location": "New Location",
  "timezone": "Europe/London"
}
```

## QSO Endpoints

### Создание QSO
```http
POST /qso
Content-Type: application/json
Authorization: Bearer <token>

{
  "callsign": "UA3ABC",
  "time_on": "2024-01-15T14:30:00Z",
  "frequency": "14.250",
  "band": "20m",
  "mode": "SSB",
  "rst_sent": "59",
  "rst_received": "59",
  "grid_square": "KO85"
}
```

**Response:**
```json
{
  "qso": {
    "id": "uuid",
    "user_id": "uuid",
    "callsign": "UA3ABC",
    "time_on": "2024-01-15T14:30:00Z",
    "time_off": "2024-01-15T14:35:00Z",
    "frequency": "14.250",
    "band": "20m",
    "mode": "SSB",
    "rst_sent": "59",
    "rst_received": "59",
    "grid_square": "KO85",
    "created_at": "2024-01-15T14:30:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  }
}
```

### Получение списка QSO
```http
GET /qso?page=1&limit=50&band=20m&mode=SSB
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` - номер страницы (default: 1)
- `limit` - количество записей (default: 50)
- `band` - фильтр по диапазону
- `mode` - фильтр по режиму
- `callsign` - фильтр по позывному
- `date_from` - фильтр по дате от
- `date_to` - фильтр по дате до

**Response:**
```json
{
  "qsos": [
    {
      "id": "uuid",
      "callsign": "UA3ABC",
      "time_on": "2024-01-15T14:30:00Z",
      "band": "20m",
      "mode": "SSB",
      "rst_sent": "59",
      "rst_received": "59",
      "grid_square": "KO85"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1250,
    "has_next": true
  }
}
```

### Получение QSO по ID
```http
GET /qso/{id}
Authorization: Bearer <token>
```

### Обновление QSO
```http
PUT /qso/{id}
Content-Type: application/json
Authorization: Bearer <token>

{
  "rst_sent": "59",
  "notes": "Good signal today"
}
```

### Удаление QSO
```http
DELETE /qso/{id}
Authorization: Bearer <token>
```

### Импорт ADIF
```http
POST /qso/import
Content-Type: multipart/form-data
Authorization: Bearer <token>

file: <adif_file>
```

**Response:**
```json
{
  "imported_count": 25,
  "error_count": 2,
  "errors": [
    "Line 15: Invalid frequency format",
    "Line 32: Missing required field 'time_on'"
  ]
}
```

### Экспорт ADIF
```http
GET /qso/export?date_from=2024-01-01&date_to=2024-01-31
Authorization: Bearer <token>
```

## Station Endpoints

### Создание станции
```http
POST /stations
Content-Type: application/json
Authorization: Bearer <token>

{
  "callsign": "RD2W",
  "name": "Home Station",
  "location": "Kursk, Russia",
  "grid_square": "KO85",
  "equipment": ["IC-7300", "LDG Z-100Plus"],
  "antennas": ["Dipole 40m", "Vertical 20m"],
  "supported_modes": ["SSB", "CW", "FT8"],
  "supported_bands": ["80m", "40m", "20m", "15m"]
}
```

### Получение списка станций
```http
GET /stations
Authorization: Bearer <token>
```

### Обновление станции
```http
PUT /stations/{id}
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Updated Station Name",
  "location": "New Location"
}
```

### Удаление станции
```http
DELETE /stations/{id}
Authorization: Bearer <token>
```

## Analytics Endpoints

### Общая статистика
```http
GET /analytics/statistics
Authorization: Bearer <token>
```

**Response:**
```json
{
  "total_qso": 1250,
  "qso_by_band": {
    "20m": 450,
    "40m": 320,
    "80m": 180
  },
  "qso_by_mode": {
    "SSB": 800,
    "CW": 300,
    "FT8": 150
  },
  "last_30_days": [25, 30, 28, ...],
  "unique_countries": 45,
  "unique_grids": 120
}
```

### Расширенная статистика
```http
GET /analytics/statistics/extended?period=30d&date_from=2024-01-01&date_to=2024-01-31
Authorization: Bearer <token>
```

**Query Parameters:**
- `period` - период (1d, 7d, 30d, 1y)
- `date_from` - начальная дата
- `date_to` - конечная дата

**Response:**
```json
{
  "basic_stats": {
    "total_qso": 1250,
    "qso_by_band": {
      "20m": 450,
      "40m": 320,
      "80m": 180
    },
    "qso_by_mode": {
      "SSB": 800,
      "CW": 300,
      "FT8": 150
    },
    "last_30_days": [25, 30, 28, ...],
    "unique_countries": 45,
    "unique_grids": 120
  },
  "growth": {
    "total_qso_growth": 15.5,
    "new_countries": 3,
    "new_grids": 8,
    "active_days": 25
  },
  "continents": [
    {
      "continent": "EU",
      "count": 650,
      "percentage": 52.0
    },
    {
      "continent": "AS",
      "count": 300,
      "percentage": 24.0
    }
  ],
  "period": "30d"
}
```

### Карта активности
```http
GET /analytics/activity?period=7d
Authorization: Bearer <token>
```

**Query Parameters:**
- `period` - период (1d, 7d, 30d, 1y)

**Response:**
```json
{
  "activity": [
    {
      "hour": "14:00",
      "band": "20m",
      "count": 5
    },
    {
      "hour": "15:00",
      "band": "40m",
      "count": 3
    }
  ],
  "period": "7d"
}
```

### Генерация отчетов
```http
POST /analytics/reports
Content-Type: application/json
Authorization: Bearer <token>

{
  "report_type": "dxcc",
  "format": "pdf",
  "date_from": "2024-01-01",
  "date_to": "2024-01-31"
}
```

**Response:**
```json
{
  "report_data": "base64_encoded_pdf",
  "report_type": "dxcc",
  "filename": "dxcc_report_202401.pdf",
  "size": 15420
}
```

## gRPC API

### Auth Service
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc UpdateProfile(UpdateProfileRequest) returns (UpdateProfileResponse);
  rpc GetProfile(GetProfileRequest) returns (GetProfileResponse);
}
```

### QSO Service
```protobuf
service QSOService {
  rpc CreateQSO(CreateQSORequest) returns (CreateQSOResponse);
  rpc GetQSO(GetQSORequest) returns (GetQSOResponse);
  rpc UpdateQSO(UpdateQSORequest) returns (UpdateQSOResponse);
  rpc DeleteQSO(DeleteQSORequest) returns (google.protobuf.Empty);
  rpc ListQSOs(ListQSOsRequest) returns (stream QSO);
  rpc SearchQSOs(SearchQSOsRequest) returns (stream QSO);
  rpc ImportADIF(ImportADIFRequest) returns (ImportADIFResponse);
  rpc ExportADIF(ExportADIFRequest) returns (ExportADIFResponse);
}
```

### Analytics Service
```protobuf
service AnalyticsService {
  rpc GetStatistics(GetStatisticsRequest) returns (GetStatisticsResponse);
  rpc GetExtendedStatistics(GetExtendedStatisticsRequest) returns (GetExtendedStatisticsResponse);
  rpc GetActivityHeatmap(GetActivityHeatmapRequest) returns (GetActivityHeatmapResponse);
  rpc GenerateReports(GenerateReportsRequest) returns (GenerateReportsResponse);
}
```

### Health Service
```protobuf
service Health {
  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
  rpc Watch(HealthCheckRequest) returns (stream HealthCheckResponse);
}
```

## Коды ошибок

| Code | Description |
|------|-------------|
| 400 | Bad Request - неверные параметры |
| 401 | Unauthorized - требуется аутентификация |
| 403 | Forbidden - недостаточно прав |
| 404 | Not Found - ресурс не найден |
| 409 | Conflict - конфликт данных |
| 422 | Unprocessable Entity - ошибка валидации |
| 429 | Too Many Requests - превышен лимит запросов |
| 500 | Internal Server Error |
| 503 | Service Unavailable - сервис временно недоступен |
