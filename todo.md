🔴 - azhaxyly
🟦 - mromanul

🧠 1. Domain Layer (Core)

Models и интерфейсы — без привязки к базе или API

Описать тип PriceUpdate (exchange, pair, price, timestamp) 🟦 ✅

Описать агрегированные данные: PriceStat (avg/min/max за минуту) 🔴 ✅

Интерфейсы:

PriceRepository (InsertStat, GetLatest, GetByPeriod, etc.) 🟦 ✅

Cache (SetLatest, GetLatest) 🟦 ✅

ExchangeClient (StartStreaming() или GetUpdates(chan PriceUpdate)) 🔴 ✅

---

⚙️ 2. Application Layer (Use Cases)

Бизнес-логика, использует интерфейсы domain

Сервис агрегации по минутам (накапливает, считает avg/min/max, сбрасывает в БД) 🔴

Сервис по обработке входящих PriceUpdate (fan-out → worker pool) 🟦 ❓✅

Сервис переключения режимов (Live/Test) 🟦 ✅

REST-сервис: GetLatestPrice(symbol), GetMax(symbol, duration) и т.п. 🔴

---

🔌 3. Adapters
💾 a) PostgreSQL (Storage Adapter) 🔴

Реализовать PriceRepository через pgx или database/sql ✅

Создать таблицу + миграцию для хранения агрегаций

Вставка батчами ✅

Запросы min/max/avg за период ✅

🧠 b) Redis (Cache Adapter) 🟦✅

Реализовать Cache (key: latest:EX:PAIR, value: цена+время)

Добавить TTL и очистку старого

Fallback → если Redis упал, не тормозить обработку

🌍 c) Exchange Clients

TCP-клиент, читающий поток с 40101 / 40102 / 40103 🔴

Автоматическое переподключение при обрыве 🔴

Генератор (Test Mode): отправляет фейковые данные по каналу 🟦

---

🌐 4. HTTP API (Web Adapter)

Реализовать хендлеры для: 

Конечные точки API 🟦🔴
API рыночных данных

GET /prices/latest/{symbol} — получить последнюю цену для заданного символа.

GET /prices/latest/{exchange}/{symbol} — получить последнюю цену для заданного символа с определенной биржи.

GET /prices/highest/{symbol} — получить самую высокую цену за период.

GET /prices/highest/{exchange}/{symbol} — получить самую высокую цену за период с определенной биржи.

GET /prices/highest/{symbol}?period={duration} — получить самую высокую цену за последний {duration} (например, последние 1 с, 3 с, 5 с, 10 с, 30 с, 1 м, 3 м, 5 м).

GET /prices/highest/{exchange}/{symbol}?period={duration} — Получить самую высокую цену за последний {duration} с определенной биржи.

GET /prices/lowest/{symbol} — Получить самую низкую цену за период.

GET /prices/lowest/{exchange}/{symbol} — Получить самую низкую цену за период с определенной биржи.

GET /prices/lowest/{symbol}?period={duration} — Получить самую низкую цену за последний {duration}.

GET /prices/lowest/{exchange}/{symbol}?period={duration} — Получить самую низкую цену за последний {duration} с определенной биржи.

GET /prices/average/{symbol} — Получить среднюю цену за период.

GET /prices/average/{exchange}/{symbol} — Получить среднюю цену за период с определенной биржи.

GET /prices/average/{exchange}/{symbol}?period={duration} — получить среднюю цену за последние {duration} с определенной биржи

API Data Mode

POST /mode/test — переключиться в тестовый режим (использовать сгенерированные данные).

POST /mode/live — переключиться в режим реального времени (извлечь данные из предоставленных программ).

Состояние системы

GET /health — возвращает состояние системы (например, соединения, доступность Redis). 

Валидация входных параметров

Ответы в формате JSON

---

🔁 5. Concurrency: Fan-In / Fan-Out / Worker Pool 🟦

Fan-In: собрать данные с 3 источников в общий chan PriceUpdate

Fan-Out: передать обновления в worker-пул

Worker Pool: обрабатывать и класть в Redis и PostgreSQL

Убедиться, что не блокируется

---

🚦 6. Graceful shutdown + Logging 🔴

Обработка SIGINT/SIGTERM

Закрытие соединений (pg, redis, источники)

Использование log/slog с уровнем Info, Error, Context ✅
