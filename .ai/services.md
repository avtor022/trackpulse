# TrackPulse - Сервисы и бизнес-логика

## Обзор сервисного слоя

Сервисы находятся в `internal/service/` и содержат всю бизнес-логику приложения.

**Направление зависимостей**: UI → Service → Repository

---

## LapService (lap_service.go)

### Назначение
Обработка RFID-сканирований в реальном времени, подсчёт кругов, ведение результатов заезда.

### Структура
```go
type LapService struct {
    rawScanRepo         *repository.RawScanRepository
    competitorModelRepo *repository.CompetitorModelRepository
    competitionRepo     *repository.CompetitionRepository
    participantRepo     *repository.CompetitionParticipantRepository
    
    scanChannel chan LapScan  // Буфер на 200 сканирований
    
    activeCompetition   *models.Competition
    transponderCache    map[string]string  // transponder -> competitor_model_id
    participantResults  map[string]*ParticipantLapData
    modelToParticipant  map[string]string  // competitor_model_id -> participant_id
    
    stopWorker chan bool
    batchSize    int           // 50 сканирований
    batchTimeout time.Duration // 100ms
}
```

### Ключевые методы

#### Инициализация и управление
```go
NewLapService(...) *LapService
Start()                    // Запуск фонового worker
Stop()                     // Остановка worker
SetActiveCompetition(comp) // Установка активного заезда
GetActiveCompetition() *Competition
```

#### Обработка и результаты
```go
QueueScan(scan LapScan) bool          // Поместить сканирование в очередь worker
ProcessScan(scan LapScan)             // Обработать сканирование синхронно
Flush()                               // Дождаться ранее поставленных в очередь сканирований
GetParticipantResults()               // Получить копию результатов активного заезда
GetStandings(competitionID)            // Получить таблицу участников и результатов из БД
PersistResults() error                 // Сохранить результаты активного заезда
```

### Алгоритм обработки сканирования

1. Worker получает сканирование и добавляет его в пакет для `raw_scans`.
2. При активном заезде ищет транспондер и участника в кэшах.
3. Первая отметка запускает первый круг, не увеличивая количество завершённых кругов.
4. Следующая отметка принимается, если прошло не меньше `MinLapTimeMs` с последней принятой отметки; более ранние отметки не меняют время начала круга.
5. Каждая принятая отметка обновляет агрегированный результат в `competition_laps`; сохранение времени выполняется с точностью до долей секунды.
6. Экран мониторинга периодически читает результаты из БД. `lap_history` пока не заполняется.

### Валидация кругов

Отметки до старта заезда, неизвестные транспондеры, отметки не позднее предыдущей и интервалы короче заданного в соревновании минимума не завершают круг.

### Потокобезопасность

- `activeCompetitionMu` — защита активного заезда
- `transponderMu` — защита кэша транспондеров
- `resultsMu` — защита результатов участников
- `cacheMu` — защита обратного lookup

---

## CompetitorService (competitor_service.go)

### Назначение
Управление участниками (пилотами): CRUD операции, валидация.

### Методы
```go
NewCompetitorService(repo) *CompetitorService

CreateCompetitor(dto) (*Competitor, error)
UpdateCompetitor(id, dto) (*Competitor, error)
DeleteCompetitor(id) error
GetByID(id) (*Competitor, error)
GetByNumber(number) (*Competitor, error)
GetAll() ([]*Competitor, error)
SearchByName(query string) ([]*Competitor, error)
ExistsByNumber(number, excludeID) bool
GenerateCompetitorNumber() int
```

### Бизнес-правила
- Номер участника уникален
- Автогенерация номера: следующий доступный
- Поиск по имени: нечувствителен к регистру, частичное совпадение

---

## RCModelService (rc_model_service.go)

### Назначение
Управление каталогом RC-моделей и справочниками.

### Методы
```go
NewRCModelService(modelRepo, brandRepo, scaleRepo, typeRepo) *RCModelService

// Модели
CreateRCModel(dto) (*RCModel, error)
UpdateRCModel(id, dto) (*RCModel, error)
DeleteRCModel(id) error
GetByID(id) (*RCModel, error)
GetAll() ([]*RCModel, error)
Search(query string) ([]*RCModel, error)

// Справочники
GetBrands() ([]string, error)
GetScales() ([]string, error)
GetTypes() ([]string, error)
AddBrand(name) error
AddScale(name) error
AddType(name) error
```

### Бизнес-правила
- Название модели уникально в рамках бренда
- Справочники заполняются автоматически при создании
- Нельзя удалить бренд/тип/масштаб, если есть ссылки

---

## CompetitorModelService (competitor_model_service.go)

### Назначение
Связывание участников с моделями и транспондерами.

### Методы
```go
NewCompetitorModelService(cmRepo, compRepo, modelRepo) *CompetitorModelService

CreateCompetitorModel(dto) (*CompetitorModel, error)
UpdateCompetitorModel(id, dto) (*CompetitorModel, error)
DeactivateCompetitorModel(id) error
GetByID(id) (*CompetitorModel, error)
GetByCompetitorID(competitorID) ([]*CompetitorModel, error)
GetByTransponder(transponder) (*CompetitorModel, error)
GetActiveByTransponder(transponder) (*CompetitorModel, error)
IsTransponderUnique(transponder, excludeID) bool
```

### Бизнес-правила
- Транспондер уникален в системе
- При привязке нового транспондера старые связи деактивируются
- Один участник может иметь несколько моделей (но один активный транспондер)

---

## CompetitionService (competition_service.go)

### Назначение
Управление заездами: создание, старт, финиш, статусы.

### Методы
```go
NewCompetitionService(compRepo, typeRepo, scaleRepo, trackRepo, yearRepo, seasonRepo) *CompetitionService

CreateCompetition(dto) (*Competition, error)
UpdateCompetition(id, dto) (*Competition, error)
DeleteCompetition(id) error
GetByID(id) (*Competition, error)
GetAll() ([]*Competition, error)
GetByStatus(status) ([]*Competition, error)
GetInProgress() (*Competition, error)
StartCompetition(id) error
FinishCompetition(id) error
CancelCompetition(id) error
ValidateCompetition(dto) error
```

### Статусы заездов
- `scheduled` — запланирован (можно редактировать)
- `in_progress` — идёт гонка (нельзя редактировать)
- `finished` — завершён (только чтение)
- `cancelled` — отменен

### Бизнес-правила
- Нельзя начать заезд без участников
- Нельзя изменить заезд после старта
- Пара (год, сезон) уникальна
- Лимит времени или кругов обязателен

---

## CompetitionParticipantService (competition_participant_service.go)

### Назначение
Управление участниками конкретного заезда.

### Методы
```go
NewCompetitionParticipantService(repo, cmService, compService) *CompetitionParticipantService

AddParticipant(competitionID, competitorModelID, gridPosition) error
RemoveParticipant(id) error
UpdateGridPosition(id, position) error
GetByCompetitionID(competitionID) ([]*CompetitionParticipant, error)
GetByCompetitorModelID(cmID) ([]*CompetitionParticipant, error)
IsRegistered(competitionID, competitorModelID) bool
MarkAsFinished(id, dnf, reason) error
```

### Бизнес-правила
- Один участник может быть зарегистрирован только один раз в заезде
- Grid position опционален
- Можно отметить как финишировавшего или DNF

---

## SettingsService (settings_service.go)

### Назначение
Управление настройками приложения.

### Методы
```go
NewSettingsService(repo) *SettingsService

GetLocale() (string, error)
SetLocale(locale string) error
GetLastCOMPort() (string, error)
SetLastCOMPort(port string) error
GetValue(key string) (string, error)
SetValue(key, value string) error
```

### Ключи настроек
- `app.locale` — язык интерфейса
- `hardware.last_port` — последний COM-порт
- `ui.theme` — тема оформления

---

## Паттерны использования

### Создание нового заезда
```go
// 1. Создать заезд
comp, err := competitionService.CreateCompetition(dto)

// 2. Добавить участников
for _, cmID := range competitorModelIDs {
    err := participantService.AddParticipant(comp.ID, cmID, nil)
}

// 3. Начать заезд
err = competitionService.StartCompetition(comp.ID)

// 4. Активировать LapService
lapService.SetActiveCompetition(comp)
lapService.Start()
```

### Обработка RFID сканирования
```go
// В UI или hardware модуле
scan := LapScan{
    TagValue: tag,
    Timestamp: time.Now(),
    ReaderType: "EM4095",
    COMPort: "COM3",
}

// Отправить в очередь сервиса
lapService.QueueScan(scan)
```

### Получение результатов заезда
```go
// Получить участников
participants, _ := participantService.GetByCompetitionID(compID)

// Для каждого получить результаты
for _, p := range participants {
    laps := getLapsForParticipant(p.ID)  // через repository
    // laps.NumberOfLaps, laps.BestLapTimeMs
}
```

---

## Интеграция с UI

### Обновление в реальном времени
- `MonitoringPanel` раз в 500 мс обновляет таблицу через `GetStandings()`.
- Обновление виджетов выполняется через `fyne.Do()`.
- Сортировка: число завершённых кругов по убыванию, лучший круг по возрастанию, номер спортсмена.

---

## Тестирование сервисов

### Пример теста LapService
```go
func TestLapService_ProcessScan_Success(t *testing.T) {
    // Setup
    mockRepos := createMockRepositories()
    svc := NewLapService(...)
    
    // Create test data
    comp := createTestCompetition("in_progress")
    svc.SetActiveCompetition(comp)
    
    // Process scan
    scan := LapScan{TagValue: "TEST123", Timestamp: time.Now()}
    svc.ProcessScan(scan)
    
    // Wait for processing
    time.Sleep(200 * time.Millisecond)
    
    // Assert
    verifyLapWasRecorded(t, ...)
}
```

### Моки репозиториев
```go
type MockCompetitorModelRepository struct {
    GetByTransponderFunc func(string) (*CompetitorModel, error)
}

func (m *MockCompetitorModelRepository) GetByTransponder(tag string) (*CompetitorModel, error) {
    return m.GetByTransponderFunc(tag)
}
```
