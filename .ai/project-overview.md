# TrackPulse - Проект для AI-агентов

## Обзор проекта

**TrackPulse** — это десктопное приложение на Go для управления соревнованиями по автомодельному спорту с RFID-хронометражем.

### Технологический стек
- **Язык**: Go 1.19
- **GUI**: Fyne v2.7.3 (кроссплатформенный UI)
- **База данных**: SQLite3 (с WAL режимом)
- **Локализация**: JSON-файлы (en, ru)
- **Аппаратная часть**: RFID-считыватели (EM4095), последовательные порты

### Основные возможности
- Управление участниками (пилотами)
- Каталог RC-моделей с брендами, типами, масштабами
- Привязка транспондеров к участникам и моделям
- Создание и управление заездами/соревнованиями
- Реальное время: обработка RFID-сканирований, подсчёт кругов
- Мониторинг заездов с таймером
- Логи и отладка

---

## Структура проекта

```
trackpulse/
├── cmd/
│   └── main.go              # Точка входа, DI контейнер
├── internal/
│   ├── config/              # Конфигурация приложения
│   │   └── config.go        # Загрузка config.json
│   ├── database/            # Работа с БД, схема
│   │   └── database.go      # Инициализация SQLite, миграции
│   ├── locale/              # Локализация (i18n)
│   │   ├── locale.go        # Система переводов
│   │   └── translations/    # en.json, ru.json
│   ├── models/              # Доменные модели/сущности
│   │   ├── entities.go      # Основные структуры (Competitor, RCModel, Competition...)
│   │   └── *.go             # Справочники (бренды, типы, масштабы, сезоны, годы, трассы)
│   ├── repository/          # Data Access Layer
│   │   ├── *_repository.go  # CRUD операции для каждой сущности
│   │   └── raw_scan_repository.go  # Массовая вставка сырых RFID-сканирований
│   ├── service/             # Бизнес-логика
│   │   ├── competitor_service.go         # Управление пилотами
│   │   ├── rc_model_service.go           # Управление моделями и справочниками
│   │   ├── competitor_model_service.go   # Привязка транспондеров
│   │   ├── competition_service.go        # Управление заездами
│   │   ├── competition_participant_service.go  # Участники заезда
│   │   ├── lap_service.go    # Обработка RFID, подсчёт кругов (асинхронный worker)
│   │   └── settings_service.go # Настройки приложения
│   └── ui/                  # Fyne UI компоненты
│       ├── app.go           # Главный UI компонент, навигация по вкладкам
│       ├── participant_panel.go  # Управление участниками заезда
│       ├── competitor_panel.go   # Управление пилотами
│       ├── rc_model_panel.go     # Каталог моделей
│       ├── competitor_model_panel.go  # Привязка транспондеров
│       ├── competition_panel.go  # Создание и настройка заездов
│       ├── competition_filter.go  # Фильтрация заездов
│       ├── monitoring_panel.go  # Экран гонки в реальном времени
│       ├── participant_panel.go  # Регистрация участников на заезд
│       ├── settings_panel.go # Настройки приложения
│       ├── logs_panel.go     # Просмотр логов
│       ├── port_scanner.go   # Сканирование COM-портов
│       ├── timer.go          # Таймер для заездов
│       ├── autocomplete_entry.go  # Поле с автодополнением
│       └── reference_popup.go # Всплывающие справочники
├── pkg/
│   ├── logger/              # Логгер (info/error/debug)
│   │   └── logger.go
│   └── utils/               # Утилиты
│       └── utils.go
├── tests/
│   ├── service/             # Юнит-тесты сервисов
│   │   ├── competitor_service_test.go
│   │   ├── competition_service_test.go
│   │   ├── rc_model_service_test.go
│   │   └── competitor_model_service_test.go
│   └── ui/                  # UI тесты
│       ├── monitoring_panel_test.go
│       └── reference_popup_test.go
├── config.json              # Конфигурационный файл
├── go.mod                   # Зависимости
└── README.md
```

---

## Архитектурные паттерны

### Clean Architecture (упрощённая)
```
UI → Service → Repository → Database
```

### Dependency Injection
В `main.go` создаются все зависимости и передаются в компоненты:
1. Logger → Config → Database
2. Repositories (все) → Services (все) → UI

### Модели данных (internal/models/entities.go)

| Сущность | Описание |
|----------|----------|
| `Competitor` | Участник (пилот): имя, номер, рейтинг, страна, город |
| `RCModel` | Модель: бренд, название, масштаб, тип, мотор, привод |
| `CompetitorModel` | Связь: участник + модель + транспондер |
| `Competition` | Заезд: название, тип, трасса, лимиты времени/кругов |
| `CompetitionParticipant` | Участник в конкретном заезде |
| `CompetitionLaps` | Агрегированные результаты участника |
| `LapHistory` | История каждого круга (детальная) |
| `RawScan` | Сырые RFID-сканирования |

### Ключевые сервисы

### LapService (lap_service.go)

**Назначение**: Обработка RFID-сканирований в реальном времени, подсчёт кругов, ведение результатов заезда.

**Особенности**:
- Фоновый worker с буферизацией (канал на 200 сканирований)
- Кэширование транспондеров и участников для O(1) поиска
- Пакетная запись в БД (50 сканирований или 100ms)
- Методы: `Start()`, `Stop()`, `SetActiveCompetition()`, `ProcessScan()`, `QueueScan()`
- Структуры: `LapScan`, `ParticipantLapData`

#### CompetitionService
- Управление заездами: создание, старт, финиш, статусы
- Валидация: лимиты времени, кругов, уникальность

#### CompetitorModelService
- Связывание участников с моделями и транспондерами
- Проверка дубликатов транспондеров

---

## База данных (SQLite)

### Основные таблицы
- `competitors` — участники
- `rc_models` — модели
- `competitor_models` — транспондеры
- `competitions` — заезды
- `competition_participants` — участники заезда
- `competition_laps` — результаты
- `lap_history` — история кругов
- `raw_scans` — сырые RFID-данные

### Справочники
- `rc_model_brands` — бренды моделей
- `rc_model_scales` — масштабы
- `rc_model_types` — типы моделей
- `competition_tracks` — трассы
- `competition_years` — годы
- `competition_seasons` — сезоны
- `settings` — настройки приложения

### Особенности
- WAL режим для производительности
- Foreign keys включены
- Индексы на часто используемых полях
- UUID для всех первичных ключей

---

## Конфигурация (config.json)

```json
{
  "db_path": "./trackpulse.db",
  "hardware_com_port": "COM3",
  "hardware_reader_type": "EM4095",
  "hardware_debounce_ms": 2000,
  "ui_language": "ru",
  "time_limit_minutes": 5,
  "default_track_name": ""
}
```

---

## Локализация

Файлы переводов в `internal/locale/translations/`:
- `en.json` — английский
- `ru.json` — русский

Использование в коде:
```go
import "trackpulse/internal/locale"

text := locale.Get("key_name")
// или
text := locale.T("key_name")
```

---

## UI Компоненты (Fyne)

### Главные панели
| Компонент | Назначение |
|-----------|------------|
| `App` | Главное окно, навигация по вкладкам |
| `ParticipantPanel` | Управление участниками заезда |
| `CompetitorPanel` | Управление пилотами |
| `RCModelPanel` | Каталог моделей |
| `CompetitorModelPanel` | Привязка транспондеров |
| `CompetitionPanel` | Создание и настройка заездов |
| `MonitoringPanel` | Экран гонки в реальном времени |
| `SettingsPanel` | Настройки приложения |
| `LogsPanel` | Просмотр логов |

### Вспомогательные компоненты
- `AutocompleteEntry` — поле с автодополнением
- `CompetitionFilter` — фильтрация заездов
- `PortScanner` — сканирование COM-портов
- `Timer` — таймер для заездов
- `ReferencePopup` — всплывающие справочники

---

## Команды для разработки

### Сборка
```bash
# Windows (GUI без консоли)
go build -ldflags="-s -w -H windowsgui" -o trackpulse.exe cmd\main.go

# Linux/Mac
go build -ldflags="-s -w" -o trackpulse cmd/main.go
```

### Запуск
```bash
# Без сборки
go run ./cmd/main.go

# Собранный бинарник
./trackpulse        # Linux/Mac
trackpulse.exe      # Windows
```

### Тесты
```bash
# Все тесты
go test ./... -v

# Конкретный пакет
go test ./tests/service/ -v -run TestCompetitorService

# Покрытие
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Очистка кэша
```bash
go clean --cache
go mod tidy
```

---

## Паттерны для AI-агентов

### При добавлении новой сущности
1. Добавить модель в `internal/models/`
2. Создать репозиторий в `internal/repository/`
3. Создать сервис в `internal/service/`
4. Добавить UI панель в `internal/ui/` (если нужно)
5. Обновить `main.go` для DI
6. Добавить переводы в `internal/locale/translations/`

### При изменении бизнес-логики
- Менять только слой **Service**
- Не менять UI напрямую из Repository
- Соблюдать направление зависимостей: UI → Service → Repository

### При работе с БД
- Использовать только методы репозиториев
- Не писать SQL в сервисах или UI
- Для новых запросов добавлять методы в репозитории

### При добавлении UI элементов
- Использовать существующие паттерны Fyne из других панелей
- Применять локализацию через `locale.T()`
- Следовать стилю существующих компонентов

---

## Важные замечания

1. **Все команды выполняются из корня проекта** (`/workspace`)
2. **UUID для ID**: Используется `github.com/google/uuid` для генерации ID
3. **Время**: Хранится как `time.Time`, в БД как TEXT (ISO8601)
4. **Транспондеры**: Уникальны в рамках системы (unique индекс)
5. **Статусы заездов**: `scheduled`, `in_progress`, `finished`, `cancelled`
6. **Потокобезопасность**: LapService использует mutex для кэшей

---

## Контакты и ресурсы

- GUI фреймворк: https://fyne.io/
- SQLite драйвер: https://github.com/mattn/go-sqlite3
- Документация Go: https://go.dev/doc/
