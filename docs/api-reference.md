# API Reference

## Общая информация

Базовый URL: `http://localhost:8080/api/v1`

### Аутентификация
Большинство endpoints требуют JWT аутентификации, поэтому необходимо добавлять токен в заголовок:
```
Authorization: Bearer <your_jwt_token>
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
  "location": "Kursk, Russia"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "callsign": "RD2W",
    "email": "user@example.com",
    "name": "Maxim Kru"
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

### Обновление токена
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token"
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

### Экспорт ADIF
```http
GET /qso/export?date_from=2024-01-01&date_to=2024-01-31
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
  "last_30_days": [25, 30, 28, ...]
}
```

### Карта активности
```http
GET /analytics/activity?period=7d
Authorization: Bearer <token>
```

**Query Parameters:**
- `period` - период (1d, 7d, 30d, 1y)

## gRPC API

### Auth Service
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
}
```

### QSO Service
```protobuf
service QSOService {
  rpc CreateQSO(CreateQSORequest) returns (CreateQSOResponse);
  rpc GetQSO(GetQSORequest) returns (GetQSOResponse);
  rpc UpdateQSO(UpdateQSORequest) returns (UpdateQSOResponse);
  rpc DeleteQSO(DeleteQSORequest) returns (DeleteQSOResponse);
  rpc ListQSOs(ListQSOsRequest) returns (stream QSO);
  rpc SearchQSOs(SearchQSOsRequest) returns (stream QSO);
}
```

### Analytics Service
```protobuf
service AnalyticsService {
  rpc GetStatistics(GetStatisticsRequest) returns (GetStatisticsResponse);
  rpc GetActivityHeatmap(GetActivityHeatmapRequest) returns (GetActivityHeatmapResponse);
  rpc GenerateReports(GenerateReportsRequest) returns (GenerateReportsResponse);
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
| 500 | Internal Server Error |
