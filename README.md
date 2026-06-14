# go-2024-2025-2
Итоговый проект «Разработка банковской системы» по дисциплине "Язык программирования GO". 2-й семестр 1-го курса МИФИ ИИКС РПО (2024-2025 уч. г).

## Описание проекта

Данный сервис представляет собой REST API для управления банковскими счетами, картами и кредитами. Подробное описание задания находится в файле [task.md](task.md).

## Требования

- Go 1.23+
- PostgreSQL 17+
- Установленные переменные окружения (см. далее)

## Установка и запуск

1. Клонирование репозитория:
```bash
git clone <repository-url>
cd go-2024-2025-2
```

2. Установка зависимостей:
```bash
go mod download
```

3. Настройка переменных окружения:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=bank
export JWT_SECRET=your_jwt_secret
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USER=your_email@example.com
export SMTP_PASSWORD=your_smtp_password
```

4. Запуск тестовой базы данных (при необходимости):
```bash
docker-compose up -d
```

5. Генерация PGP ключей:
```bash
go run scripts/generate_key.go
```

6. Запуск сервиса:
```bash
go run cmd/api/main.go
```

Сервис будет доступен по адресу `http://localhost:8080`

## Использование API

### Аутентификация

- **Регистрация**  
  `POST /api/register`  
  Тело запроса:
  ```json
  {
    "username": "имя_пользователя",
    "email": "email@example.com",
    "password": "пароль"
  }
  ```

- **Вход**  
  `POST /api/login`  
  Тело запроса:
  ```json
  {
    "email": "email@example.com",
    "password": "пароль"
  }
  ```
  Возвращает JWT токен для авторизованных запросов.

### Управление счетами

- **Создать счет**  
  `POST /api/accounts/create`  
  Тело запроса:
  ```json
  {
    "type": "DEBIT"
  }
  ```

- **Получить список счетов**  
  `GET /api/accounts/list`

- **Пополнить счет**  
  `POST /api/accounts/deposit`  
  Тело запроса:
  ```json
  {
    "account_id": 1,
    "amount": 1000.00
  }
  ```

- **Снять средства**  
  `POST /api/accounts/withdraw`  
  Тело запроса:
  ```json
  {
    "account_id": 1,
    "amount": 500.00
  }
  ```

- **Перевод между счетами**  
  `POST /api/accounts/transfer`  
  Тело запроса:
  ```json
  {
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 300.00
  }
  ```

### Управление картами

- **Создать виртуальную карту**  
  `POST /api/cards/create`  
  Тело запроса:
  ```json
  {
    "account_id": 1
  }
  ```

- **Получить список карт**  
  `GET /api/cards/list`

- **Получить информацию о карте**  
  `GET /api/cards/get?card_id=1`

### Кредиты

- **Оформить кредит**  
  `POST /api/credits/create`  
  Тело запроса:
  ```json
  {
    "account_id": 1,
    "amount": 10000.00,
    "term": 12,
    "rate": 15.5
  }
  ```

- **Получить список кредитов**  
  `GET /api/credits/list`

- **Получить информацию о кредите**  
  `GET /api/credits/get?id=1`

- **Получить график платежей**  
  `GET /api/credits/schedule?id=1`

- **Создать платеж**  
  `POST /api/payments/create`  
  Тело запроса:
  ```json
  {
    "credit_id": 1,
    "amount": 1000.00,
    "due_date": "2024-05-01T00:00:00Z"
  }
  ```

- **Обработать платеж**  
  `POST /api/payments/process?payment_id=1`

- **Получить платежи по кредиту**  
  `GET /api/payments/list?credit_id=1`

- **Получить ожидающие платежи**  
  `GET /api/payments/pending`

## Особенности реализации

### Безопасность
- Пароли хешируются с помощью bcrypt
- Данные карт шифруются с помощью PGP
- CVV хешируется с помощью bcrypt
- Все запросы защищены JWT-аутентификацией
- Используются параметризованные SQL-запросы

### Интеграции
- SMTP для отправки уведомлений
- SOAP API ЦБ РФ для получения ключевой ставки
- Автоматическое списание платежей по кредитам

### Логирование
Логи сохраняются в файл `app.log` и выводятся в консоль. Используется logrus с настройками:
- Уровень логирования: Info
- Формат: JSON
- Выход: файл и консоль

## Структура проекта

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   └── utils/
├── migrations/
├── scripts/
└── README.md
```

## Технические детали

- Используется PostgreSQL с расширением pgcrypto
- Реализован алгоритм Луна для генерации номеров карт
- Расчет аннуитетных платежей для кредитов
- Транзакции для обеспечения атомарности операций
- Middleware для аутентификации и логирования
