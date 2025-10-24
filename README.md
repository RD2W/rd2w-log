# Ham Radio QSO Journal

[![Go Version](https://img.shields.io/badge/Go-1.25.3+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GitHub Issues](https://img.shields.io/github/issues/rd2w/rd2w-log)](https://github.com/rd2w/rd2w-log/issues)

Современная микросервисная система для ведения журнала радиосвязей (QSO) с аналитикой и поддержкой мобильных клиентов.

## 🚀 Быстрый старт

```bash
# Клонирование репозитория
git clone https://github.com/rd2w/rd2w-log.git
cd rd2w-log

# Запуск в Docker
docker-compose up -d

# Проверка работы
curl http://localhost:8080/api/v1/health
```

## 📖 Документация

- [📋 Техническое задание](docs/specification.md) - полное описание проекта и архитектуры
- [🔗 API Reference](docs/api-reference.md) - документация по API
- [🚀 Deployment Guide](docs/deployment.md) - инструкции по развертыванию

## 🏗️ Архитектура

Проект построен на основе микросервисной архитектуры с использованием:

- **Go 1.25.3+** с чистой архитектурой
- **gRPC** + **Protocol Buffers** для межсервисной коммуникации
- **PostgreSQL** + **Redis** для хранения и кэширования
- **JWT** для аутентификации
- **Docker** для контейнеризации

## 🛠️ Основные сервисы

- **Auth Service** - аутентификация и управление пользователями
- **QSO Service** - управление журналом радиосвязей
- **Analytics Service** - статистика и аналитика
- **API Gateway** - REST/JSON интерфейс

## 📊 Функциональность

- ✅ Ведение журнала QSO (CRUD операции)
- ✅ Импорт/экспорт в формате ADIF
- ✅ Поиск и фильтрация записей
- ✅ Базовая статистика и аналитика
- ✅ Аутентификация и авторизация
- ✅ Управление станциями

## 🔮 Планы развития

- [ ] Интеграция с внешними API (QRZ.com, HamQTH)
- [ ] Расширенная аналитика (DXCC, awards)
- [ ] Real-time уведомления
- [ ] Мобильное приложение

## 🤝 Участие в разработке

Мы приветствуем вклад в проект! Пожалуйста, ознакомьтесь с [руководством по Contributing](CONTRIBUTING.md) перед началом работы.

## 📄 Лицензия

Этот проект лицензирован under the Apache 2.0 License - смотрите файл [LICENSE](LICENSE) для деталей.

---

**Разработано радиолюбителем для радиолюбителей** 🎯
