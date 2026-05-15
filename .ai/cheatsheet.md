# TrackPulse - Шпаргалка для AI-агентов

## Быстрый старт

### 1. Запуск проекта
```bash
cd /workspace
go run ./cmd/main.go
```

### 2. Сборка
```bash
# Windows
go build -ldflags="-s -w -H windowsgui" -o trackpulse.exe cmd\main.go

# Linux/Mac
go build -ldflags="-s -w" -o trackpulse cmd/main.go
```

### 3. Тесты
```bash
go test ./... -v
```

---

## Структура проекта (кратко)

```
trackpulse/
├── cmd/main.go              # Точка входа, DI
├── internal/
│   ├── config/              # Конфигурация
│   ├── database/            # БД (SQLite)
│   ├── locale/              # Переводы (en, ru)
│   ├── models/              # Сущности
│   ├── repository/          # CRUD операции
│   ├── service/             # Бизнес-логика
│   └── ui/                  # Fyne UI
├── pkg/
│   ├── logger/              # Логгер
│   └── utils/               # Утилиты
└── tests/                   # Тесты
```

---

## Архитектура

```
UI → Service → Repository → Database
```

**Правило**: Никогда не нарушайте направление зависимостей!

---

## Основные сущности

| Сущность | Файл | Описание |
|----------|------|----------|
| Competitor | `models/entities.go` | Пилот (номер, имя, рейтинг) |
| RCModel | `models/entities.go` | Модель (бренд, тип, масштаб) |
| CompetitorModel | `models/entities.go` | Связь: пилот + модель + транспондер |
| Competition | `models/entities.go` | Заезд (название, статус, лимиты) |
| CompetitionParticipant | `models/entities.go` | Участник в заезде |
| CompetitionLaps | `models/entities.go` | Результаты участника |
| LapHistory | `models/entities.go` | История кругов |

---

## Сервисы (business logic)

| Сервис | Файл | Назначение |
|--------|------|------------|
| CompetitorService | `competitor_service.go` | CRUD пилотов |
| RCModelService | `rc_model_service.go` | CRUD моделей + справочники |
| CompetitorModelService | `competitor_model_service.go` | Привязка транспондеров |
| CompetitionService | `competition_service.go` | Управление заездами |
| CompetitionParticipantService | `competition_participant_service.go` | Участники заезда |
| LapService | `lap_service.go` | **RFID сканирования, подсчёт кругов** |
| SettingsService | `settings_service.go` | Настройки приложения |

---

## Ключевые файлы

### Конфигурация
- `config.json` — настройки (БД, COM-порт, язык)
- `internal/config/config.go` — загрузка конфигурации

### База данных
- `internal/database/database.go` — схема БД, инициализация

### Локализация
- `internal/locale/locale.go` — система переводов
- `internal/locale/translations/en.json` — английский
- `internal/locale/translations/ru.json` — русский

Использование:
```go
text := locale.T("key_name")
```

### Главный файл
- `cmd/main.go` — создание всех зависимостей (DI контейнер)

---

## Паттерны кода

### Создание новой сущности

1. **Модель** (`internal/models/`):
```go
type NewEntity struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

2. **Репозиторий** (`internal/repository/`):
```go
type NewEntityRepository struct {
    db *sql.DB
}

func NewNewEntityRepository(db *sql.DB) *NewEntityRepository {
    return &NewEntityRepository{db}
}

func (r *NewEntityRepository) Create(entity *models.NewEntity) error {
    _, err := r.db.Exec(
        "INSERT INTO new_entities (...) VALUES (...)",
        entity.ID, entity.Name, ...,
    )
    return err
}
```

3. **Сервис** (`internal/service/`):
```go
type NewEntityService struct {
    repo *repository.NewEntityRepository
}

func NewNewEntityService(repo *repository.NewEntityRepository) *NewEntityService {
    return &NewEntityService{repo}
}

func (s *NewEntityService) Create(dto *NewEntityDTO) (*models.NewEntity, error) {
    // Валидация
    if dto.Name == "" {
        return nil, errors.New("name is required")
    }
    
    // Создание
    entity := &models.NewEntity{
        ID: uuid.New().String(),
        Name: dto.Name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    err := s.repo.Create(entity)
    return entity, err
}
```

4. **UI панель** (`internal/ui/`):
```go
type NewEntityPanel struct {
    service *service.NewEntityService
    window  fyne.Window
    table   *widget.Table
}

func NewNewEntityPanel(service *service.NewEntityService) *NewEntityPanel {
    p := &NewEntityPanel{service: service}
    p.table = widget.NewTable(...)
    return p
}
```

5. **Обновить main.go**:
```go
newEntityRepo := repository.NewNewEntityRepository(db.DB)
newEntityService := service.NewNewEntityService(newEntityRepo)
uiApp := ui.NewApp(..., newEntityService, ...)
```

6. **Добавить переводы** (`internal/locale/translations/*.json`):
```json
{
  "new_entity.list.title": "New Entities",
  "new_entity.add": "Add New Entity"
}
```

---

### Работа с базой данных

**Правильно** (через репозиторий):
```go
// В сервисе
competitor, err := s.repo.GetByID(id)
```

**Неправильно** (SQL в сервисе):
```go
// НЕ ДЕЛАЙТЕ ТАК!
rows, err := s.db.Query("SELECT * FROM competitors WHERE id = ?", id)
```

---

### Обработка ошибок

```go
func (s *Service) DoSomething() error {
    result, err := s.repo.GetData()
    if err != nil {
        if strings.Contains(err.Error(), "UNIQUE constraint") {
            return errors.New("record already exists")
        }
        return fmt.Errorf("failed to get data: %w", err)
    }
    
    if result == nil {
        return errors.New("not found")
    }
    
    return nil
}
```

---

### UI: Обновление из горутины

```go
// Правильно
go func() {
    data, err := s.service.GetData()
    
    fyne.Do(func() {
        if err != nil {
            dialog.ShowError(err, window)
            return
        }
        table.Refresh()
    })
}()

// Неправильно (вызовет панику!)
go func() {
    data, _ := s.service.GetData()
    table.Refresh()  // Паника!
}()
```

---

### UI: Локализация

```go
// Всегда используйте locale.T()
label := widget.NewLabel(locale.T("my.key"))
button := widget.NewButton(locale.T("action.save"), handler)
entry.SetPlaceHolder(locale.T("field.placeholder"))

// НЕ используйте хардкод текстов!
label := widget.NewLabel("Save")  // ПЛОХО!
```

---

## Команды Go

### Модули
```bash
go mod tidy           # Обновить зависимости
go mod download       # Скачать зависимости
go mod verify         # Проверить целостность
```

### Сборка
```bash
go build -o app       # Обычная сборка
go build -ldflags="-s -w" -o app  # Без отладки (меньше размер)
go build -ldflags="-s -w -H windowsgui" -o app.exe  # Windows GUI
```

### Тесты
```bash
go test ./...                    # Все тесты
go test ./pkg/... -v             # Тесты пакета с выводом
go test -run TestName ./...      # Конкретный тест
go test -cover ./...             # Покрытие кода
go test -coverprofile=coverage.out
go tool cover -html=coverage.out # HTML отчёт
```

### Очистка
```bash
go clean --cache      # Очистить кэш компиляции
go clean --modcache   # Очистить кэш модулей
go clean -i ./...     # Очистить установленные пакеты
```

### Форматирование
```bash
go fmt ./...          # Форматировать код
go vet ./...          # Проверка на ошибки
goimports -w .        # Автоимпорты (нужен goimports)
```

---

## SQLite команды

### Подключение
```bash
sqlite3 trackpulse.db
```

### Полезные запросы
```sql
-- Показать все таблицы
.tables

-- Схема таблицы
.schema competitors

-- Последние заезды
SELECT * FROM competitions ORDER BY created_at DESC LIMIT 10;

-- Участники с результатами
SELECT c.full_name, cl.number_of_laps, cl.best_lap_time_ms
FROM competition_laps cl
JOIN competition_participants cp ON cl.competition_participant_id = cp.id
JOIN competitor_models cm ON cp.competitor_model_id = cm.id
JOIN competitors c ON cm.competitor_id = c.id
WHERE cp.competition_id = 'uuid-here'
ORDER BY cl.number_of_laps DESC;

-- Активные транспондеры
SELECT transponder_number, c.full_name, rm.model_name
FROM competitor_models cm
JOIN competitors c ON cm.competitor_id = c.id
JOIN rc_models rm ON cm.rc_model_id = rm.id
WHERE cm.is_active = 1;

-- Статистика по пилоту
SELECT COUNT(*) as races, SUM(cl.number_of_laps) as total_laps
FROM competition_participants cp
JOIN competitor_models cm ON cp.competitor_model_id = cm.id
LEFT JOIN competition_laps cl ON cp.id = cl.competition_participant_id
WHERE cm.competitor_id = 'uuid-here';
```

---

## Отладка

### Логи
Файлы логов находятся в папке `logs/`:
- `info_YYYY-MM-DD.log` — информационные сообщения
- `error_YYYY-MM-DD.log` — ошибки
- `debug_YYYY-MM-DD.log` — отладочная информация

### Включение debug режима
В коде используется `log.Debug()`, логи пишутся всегда.

### Просмотр логов в реальном времени
```bash
# Linux/Mac
tail -f logs/info_$(date +%Y-%m-%d).log

# Windows PowerShell
Get-Content logs\info_$(Get-Date -Format yyyy-MM-dd).log -Wait -Tail 50
```

---

## Частые задачи

### Добавить новый справочник

1. Создать таблицу в `database.go`:
```sql
CREATE TABLE IF NOT EXISTS new_dicts (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
```

2. Добавить модель в `models/`:
```go
type NewDict struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

3. Создать репозиторий в `repository/`
4. Добавить методы в существующий сервис или создать новый
5. Добавить в UI (если нужно)

---

### Изменить поле сущности

1. Обновить struct в `models/entities.go`
2. Добавить миграцию в `database.go` (ALTER TABLE)
3. Обновить репозиторий (SQL запросы)
4. Обновить сервис (валидация, DTO)
5. Обновить UI (формы, таблицы)
6. Обновить переводы

---

### Добавить новую локализацию

1. Создать файл `internal/locale/translations/de.json`
2. Скопировать структуру из `en.json`
3. Перевести все значения
4. Добавить в `SupportedLocales` в `locale.go`:
```go
var SupportedLocales = map[string]string{
    "en": "English",
    "ru": "Русский",
    "de": "Deutsch",  // Новый
}
```

---

## Советы

### ✅ Делайте
- Следуйте архитектуре UI → Service → Repository
- Используйте `locale.T()` для всех текстов
- Обрабатывайте ошибки явно
- Пишите тесты для сервисов
- Обновляйте документацию при изменении кода
- Используйте UUID для всех ID
- Кэшируйте данные в LapService

### ❌ Не делайте
- Не пишите SQL в сервисах или UI
- Не обновляйте UI из горутин без `fyne.Do()`
- Не хардкодьте тексты
- Не игнорируйте ошибки
- Не нарушайте направление зависимостей
- Не храните пароли/секреты в коде

---

## Ресурсы

- **Fyne Docs**: https://fyne.io/docs/
- **Go Docs**: https://go.dev/doc/
- **SQLite**: https://www.sqlite.org/docs.html
- **UUID**: https://github.com/google/uuid

---

## Контакты в коде

- Главный разработчик: см. `git log`
- Вопросы по UI: см. `internal/ui/app.go`
- Вопросы по БД: см. `internal/database/database.go`
- Вопросы по RFID: см. `internal/service/lap_service.go`
