# Agents Skills — TrackPulse Project

## О проекте

**TrackPulse** — RFID-приложение на Go с графическим интерфейсом (Fyne), SQLite базой данных и поддержкой локализации.

### Технологический стек

- **Язык:** Go 1.19+
- **GUI:** Fyne v2.7.3
- **База данных:** SQLite (go-sqlite3)
- **Локализация:** nicksnyder/go-i18n/v2
- **UUID:** google/uuid
- **Тестирование:** testify

---

## Структура проекта

```
trackpulse/
├── cmd/                  # Точка входа приложения
│   └── main.go
├── internal/             # Внутренние пакеты (private)
│   ├── config/           # Конфигурация приложения
│   ├── database/         # Работа с БД
│   ├── locale/           # Локализация (translations/)
│   ├── models/           # Модели данных
│   ├── repository/       # Слой доступа к данным
│   ├── service/          # Бизнес-логика
│   └── ui/               # GUI компоненты (Fyne)
├── pkg/                  # Публичные пакеты
│   ├── logger/           # Логирование
│   └── utils/            # Утилиты
├── tests/                # Тесты
│   ├── service/          # Тесты сервисов
│   └── ui/               # Тесты UI
├── config.json           # Конфигурационный файл
├── go.mod                # Зависимости Go
└── README.md             # Документация
```

---

## Навыки для работы с проектом

### 1. Сборка и запуск

#### Запуск без сборки
```bash
go clean --cache
go mod tidy
go run ./cmd/main.go
```

#### Сборка под Linux
```bash
go build -ldflags="-s -w" -o trackpulse cmd/main.go
chmod +x trackpulse
./trackpulse
```

#### Сборка под Windows
```bash
go build -ldflags="-s -w -H windowsgui" -o trackpulse.exe cmd\main.go
trackpulse.exe
```

#### Сборка под MacOS
```bash
go build -ldflags="-s -w" -o trackpulse cmd/main.go
```

---

### 2. Тестирование

#### Запуск всех тестов
```bash
go test ./... -v
```

#### Запуск тестов сервиса
```bash
go test ./tests/service/ -v -run TestCompetitorService
```

#### Запуск конкретного теста
```bash
go test ./tests/service/ -v -run TestCompetitorService_CreateCompetitor_Success
```

#### Отчёт о покрытии кода
```bash
go test ./tests/service/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

### 3. Работа с конфигурацией

Файл: `config.json`

```json
{
  "db_path": "./trackpulse.db",
  "hardware_com_port": "",
  "hardware_reader_type": "EM4095",
  "hardware_debounce_ms": 2000,
  "ui_language": "en",
  "time_limit_minutes": 5,
  "default_track_name": ""
}
```

**Параметры:**
- `db_path` — путь к SQLite базе данных
- `hardware_com_port` — COM-порт для RFID-считывателя
- `hardware_reader_type` — тип считывателя (EM4095)
- `hardware_debounce_ms` — задержка между чтениями (мс)
- `ui_language` — язык интерфейса (en/ru)
- `time_limit_minutes` — лимит времени (минуты)
- `default_track_name` — имя трека по умолчанию

---

### 4. Локализация

Файлы переводов: `internal/locale/translations/`

- `en.json` — английская локаль
- `ru.json` — русская локаль

**Добавление нового перевода:**
1. Открыть соответствующий JSON-файл
2. Добавить ключ-значение в формате: `"key": "translated text"`
3. Сохранить файл

---

### 5. Работа с базой данных

- **Driver:** github.com/mattn/go-sqlite3
- **Файл БД:** trackpulse.db (по умолчанию)

**Важно:** При изменениях схемы БД обновлять миграции в `internal/database/`.

---

### 6. Модульная структура

#### internal/config
Загрузка и управление конфигурацией приложения.

#### internal/database
Подключение к SQLite, выполнение миграций, транзакции.

#### internal/models
Структуры данных:
- Competitor (участник)
- Track (трек)
- Result (результат заезда)

#### internal/repository
CRUD-операции для моделей. Паттерн Repository.

#### internal/service
Бизнес-логика приложения. Валидация, обработка данных.

#### internal/ui
GUI на Fyne:
- Главное окно
- Формы ввода
- Таблицы результатов
- Настройки

#### pkg/logger
Логирование событий приложения.

#### pkg/utils
Вспомогательные функции (парсинг времени, форматирование).

---

### 7. Зависимости

Основные зависимости (`go.mod`):

| Пакет | Назначение |
|-------|-----------|
| fyne.io/fyne/v2 | GUI фреймворк |
| github.com/mattn/go-sqlite3 | SQLite драйвер |
| github.com/nicksnyder/go-i18n/v2 | Локализация |
| github.com/google/uuid | Генерация UUID |
| github.com/stretchr/testify | Тестирование |

**Обновление зависимостей:**
```bash
go get -u
go mod tidy
```

---

## Рекомендации для агентов

1. **Всегда работайте из корня проекта** (`/workspace`)
2. **Проверяйте модули** перед запуском: `go mod tidy`
3. **Запускайте тесты** после изменений: `go test ./...`
4. **Соблюдайте структуру** — не создавайте файлы вне `/internal`, `/pkg`, `/cmd`, `/tests`
5. **Локализация** — все пользовательские строки должны быть в файлах переводов
6. **Конфигурация** — используйте `config.json` для настроек, не хардкодьте значения
7. **Документация** — при изменении кода обновляйте соответствующую документацию (README.md, agents.md, комментарии в коде)
