# ExchangeRateService

GRPC сервис для получения курса валют с биржи Garantex с сохранением в PostgreSQL базу данных.

## 📋 Описание

Сервис реализует следующий функционал:
- Получение курса (ask и bid цены) с биржи Garantex
- Сохранение курса с меткой времени в PostgreSQL при каждом запросе
- GRPC API для получения курсов валют
- Healthcheck для проверки работоспособности
- Автоматические миграции БД

## 🚀 Быстрый старт

### Запуск через Docker Compose

```bash
# Клонируем репозиторий
git clone https://github.com/KVSH-user/ExchangeRateService.git
cd ExchangeRateService

# Запускаем все сервисы
docker-compose up -d

# Проверяем логи
docker-compose logs -f exchange-rate-service
```

## 📡 API

### GRPC Methods

#### GetRates
Получает курс валют с биржи Garantex и сохраняет в БД.

**Request:**
```protobuf
message GetRatesRequest {
  string market = 1;  // Рынок (например: "usdtrub", "btcrub")
}
```

**Response:**
```protobuf
message GetRatesResponse {
  int64 ts = 1;                        // Timestamp получения курса
  google.type.Decimal ask_price = 2;   // Цена продажи
  google.type.Decimal bid_price = 3;   // Цена покупки
}
```

**Пример вызова:**
```bash
grpcurl -plaintext -d '{"market":"usdtrub"}' localhost:9049 exchangerateservice.ExchangeRateService/GetRates
```

#### HealthCheck
Проверка работоспособности сервиса.

**Response:**
```protobuf
message HealthCheckResponse {
  string status = 1;  // "OK" если сервис работает
}
```

**Пример вызова:**
```bash
grpcurl -plaintext localhost:9049 exchangerateservice.ExchangeRateService/HealthCheck
```

## 🐳 Docker

### Docker Compose

Сервис включает полную Docker Compose конфигурацию:

```yaml
services:
  exchange-rate-service:    # Основное приложение
  postgresql:              # База данных PostgreSQL
```

**Порты:**
- `9049` - GRPC сервер
- `5432` - PostgreSQL (для внешнего доступа)


## 📄 Лицензия

Этот проект распространяется под лицензией MIT. См. файл [LICENSE](LICENSE) для деталей.
