# Database Schema

## Таблицы базы данных

### users
Основная таблица пользователей.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID пользователя |
| email | VARCHAR(255) UNIQUE | Email пользователя |
| has_secret_access | BOOLEAN | Есть ли секретный доступ |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления (soft delete) |

---

### otps
OTP коды для авторизации.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID записи |
| email | VARCHAR(255) | Email получателя |
| code | VARCHAR(10) | OTP код |
| expires_at | TIMESTAMP | Время истечения |
| verified | BOOLEAN | Проверен ли код |
| created_at | TIMESTAMP | Дата создания |

---

### categories
Категории мини-приложений.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID категории |
| name | VARCHAR(255) UNIQUE | Название категории |
| slug | VARCHAR(255) UNIQUE | URL-slug |
| description | TEXT | Описание |
| icon | VARCHAR(255) | Иконка |
| order | INTEGER | Порядок сортировки |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### mini_apps
Мини-приложения.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID приложения |
| title | VARCHAR(255) | Название |
| subtitle | VARCHAR(255) | Подзаголовок |
| icon | VARCHAR(255) | Иконка |
| category_id | INTEGER FK | ID категории |
| is_secret | BOOLEAN | Секретное приложение? |
| users_count | INTEGER | Количество пользователей |
| status | VARCHAR(50) | Статус (verified, new) |
| description | TEXT | Описание |
| features | JSONB | Список возможностей |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### cards
Информационные карты (Crypto.com, Binance, etc.).

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID карты |
| name | VARCHAR(255) | Название |
| icon | VARCHAR(255) | Иконка/логотип |
| users_count | INTEGER | Количество пользователей |
| description | TEXT | Описание |
| features | JSONB | Список возможностей |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### solana_cards
Виртуальные карты пользователей.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID карты |
| user_id | INTEGER FK | ID пользователя |
| card_number | VARCHAR(19) | Номер карты |
| cardholder_name | VARCHAR(255) | Имя владельца |
| expiry_month | VARCHAR(2) | Месяц истечения |
| expiry_year | VARCHAR(4) | Год истечения |
| cvv | VARCHAR(3) | CVV код |
| balance | DECIMAL(10,2) | Баланс |
| currency | VARCHAR(3) | Валюта (USD) |
| is_active | BOOLEAN | Активна ли карта |
| card_type | VARCHAR(20) | Тип (visa, mastercard) |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### transactions
История транзакций карт.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID транзакции |
| card_id | INTEGER FK | ID карты |
| type | VARCHAR(50) | Тип (purchase, topup, transfer) |
| amount | DECIMAL(10,2) | Сумма |
| currency | VARCHAR(3) | Валюта |
| description | TEXT | Описание |
| merchant | VARCHAR(255) | Продавец |
| status | VARCHAR(50) | Статус (completed, pending, failed) |
| crypto_type | VARCHAR(10) | Тип крипты (SOL, USDC, etc.) |
| crypto_amount | DECIMAL(20,8) | Количество крипты |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### secret_numbers
Виртуальные номера для Secret Login.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID номера |
| number | VARCHAR(20) UNIQUE | Номер телефона |
| country_code | VARCHAR(2) | Код страны |
| price | DECIMAL(10,2) | Цена |
| currency | VARCHAR(3) | Валюта |
| is_available | BOOLEAN | Доступен ли номер |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

### secret_accesses
Активированный секретный доступ.

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | ID доступа |
| user_id | INTEGER FK | ID пользователя |
| secret_number_id | INTEGER FK | ID номера |
| activated_at | TIMESTAMP | Дата активации |
| expires_at | TIMESTAMP | Дата истечения (nullable) |
| is_active | BOOLEAN | Активен ли доступ |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления |
| deleted_at | TIMESTAMP | Дата удаления |

---

## Связи между таблицами

- `mini_apps.category_id` → `categories.id`
- `solana_cards.user_id` → `users.id`
- `transactions.card_id` → `solana_cards.id`
- `secret_accesses.user_id` → `users.id`
- `secret_accesses.secret_number_id` → `secret_numbers.id`

---

## Индексы

GORM автоматически создаст индексы для:
- Всех PRIMARY KEY
- Всех FOREIGN KEY
- Полей с тегом `unique`
- Полей `deleted_at` (для soft delete)

---

## Миграции

База данных автоматически мигрируется при запуске сервера через `gorm.AutoMigrate()`.

Для ручного контроля миграций можно использовать инструменты:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
