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

#### Обработка сканирований
```go
ProcessScan(scan LapScan)          // Добавить сканирование в канал
processingWorker()                 // Фоновый обработчик
processSingleScan(scan)            // Обработка одного сканирования
identifyCompetitor(tagValue)       // Поиск по транспондеру
```

#### Работа с кругами
```go
calculateLap(participantID, scanTime)  // Расчёт времени круга
updateLapData(participantID, lapTime)  // Обновление данных
isValidLap(lapTime, participant) bool  // Валидация круга
```

#### Сохранение в БД
```go
saveBatch()                        // Пакетная запись
saveRawScans(scans)                // Сохранить сырые данные
saveLapHistory(lap)                // Сохранить историю круга
updateCompetitionLaps(participant) // Обновить агрегированные результаты
```

#### Кэширование
```go
loadParticipantCache(competitionID)  // Загрузить участников в память
rebuildCaches()                      // Перестроить все кэши
getParticipantByTransponder(tag) string
```

### Алгоритм обработки сканирования

1. **Получение сканирования** из канала
2. **Поиск транспондера** в кэше → получение competitor_model_id
3. **Поиск участника** в modelToParticipant → получение participant_id
4. **Проверка статуса заезда** (должен быть "in_progress")
5. **Расчёт времени круга**:
   - Если это первый круг → сохранение времени старта
   - Иначе → вычисление разницы с последним прохождением
6. **Валидация круга**:
   - Минимальное время круга (> 2 секунд)
   - Максимальное время (< время заезда)
   - Не дубликат последнего сканирования
7. **Сохранение**:
   - Сырое сканирование в RawScan
   - История круга в LapHistory
   - Обновление CompetitionLaps
8. **Обновление UI** (через callback или событие)

### Валидация кругов

Круг считается невалидным если:
- Время круга < 2000 мс (false positive)
- Время круга > общего времени заезда
- Это повторное сканирование того же круга
- Заезд ещё не начался или уже закончен

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

// Отправить в сервис
lapService.ProcessScan(scan)
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

### События от LapService
```go
// В monitoring_panel.go
type MonitoringPanel struct {
    lapService *service.LapService
    updateChan chan LapUpdate
}

// Callback для обновления UI
func (p *MonitoringPanel) onLapReceived(lap *LapData) {
    // Обновить таблицу результатов
    // Подсветить новый круг
    // Обновить таймеры
}
```

### Обновление в реальном времени
- LapService отправляет события в UI канал
- UI обновляется в главном потоке (Fyne требует)
- Используется timer для периодического обновления

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
