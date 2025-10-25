# API Examples

## Аутентификация

### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "callsign": "RD2W",
    "email": "user@example.com",
    "password": "securepassword",
    "name": "Maxim Kru",
    "location": "Kursk, Russia",
    "timezone": "Europe/Moscow"
  }'
```

### Вход в систему
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "callsign": "RD2W",
    "password": "securepassword"
  }'
```

## Работа с QSO

### Создание QSO
```bash
curl -X POST http://localhost:8080/api/v1/qso \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "callsign": "UA3ABC",
    "time_on": "2024-01-15T14:30:00Z",
    "frequency": "14.250",
    "band": "20m",
    "mode": "SSB",
    "rst_sent": "59",
    "rst_received": "59",
    "grid_square": "KO85"
  }'
```

### Получение списка QSO с фильтрацией
```bash
curl "http://localhost:8080/api/v1/qso?page=1&limit=50&band=20m&mode=SSB" \
  -H "Authorization: Bearer <your_jwt_token>"
```

## Аналитика

### Получение статистики
```bash
curl "http://localhost:8080/api/v1/analytics/statistics?date_from=2024-01-01&date_to=2024-01-31" \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Генерация отчета
```bash
curl -X POST http://localhost:8080/api/v1/analytics/reports \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "report_type": "dxcc",
    "format": "pdf",
    "date_from": "2024-01-01",
    "date_to": "2024-01-31"
  }'
```

## gRPC примеры

### Тестирование gRPC сервиса
```bash
# Проверка здоровья
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# Регистрация пользователя
grpcurl -plaintext -d '{
  "callsign": "RD2W",
  "email": "user@example.com",
  "password": "password",
  "name": "Test User",
  "location": "Test Location"
}' localhost:50051 auth.AuthService/Register
```
