# TrackPulse - Модели данных

## Основные сущности

### Competitor (Участник/Пилот)
```go
type Competitor struct {
    ID               string     // UUID
    CompetitorNumber int        // Уникальный номер пилота
    FullName         string     // ФИО
    Birthday         *time.Time // Дата рождения (опционально)
    Country          string     // Страна
    City             string     // Город
    Rating           int        // Рейтинг
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

**Индексы**: по номеру, по имени

---

### RCModel (RC Модель)
```go
type RCModel struct {
    ID        string    // UUID
    Brand     string    // Бренд (ссылка на справочник)
    ModelName string    // Название модели
    Scale     string    // Масштаб (1:10, 1:8 и т.д.)
    ModelType string    // Тип (Touring, Buggy, Truggy...)
    MotorType string    // Тип мотора
    DriveType string    // Привод (2WD, 4WD)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Индексы**: по бренду, по типу, по масштабу

---

### CompetitorModel (Связь: Пилот + Модель + Транспондер)
```go
type CompetitorModel struct {
    ID                string    // UUID
    CompetitorID      string    // Ссылка на Competitor
    RCModelID         string    // Ссылка на RCModel
    TransponderNumber string    // Номер транспондера (уникальный)
    TransponderType   string    // Тип (RFID, AMB...)
    IsActive          bool      // Активна ли связь
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

**Индексы**: по транспондеру, по пилоту, по модели  
**Важно**: `TransponderNumber` уникален в системе

---

### Competition (Заезд/Соревнование)
```go
type Competition struct {
    ID               string     // UUID
    CompetitionTitle string     // Название заезда
    CompetitionType  string     // Тип (qualifying, final...)
    ModelType        string     // Тип моделей
    ModelScale       string     // Масштаб
    TrackName        string     // Трасса
    LapCountTarget   *int       // Целевое количество кругов
    TimeLimitMinutes *int       // Лимит времени в минутах
    TimeStart        *time.Time // Время старта
    TimeFinish       *time.Time // Время финиша
    Status           string     // scheduled|in_progress|finished|cancelled
    CompetitionYear  *int       // Год
    Season           string     // Сезон
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

**Индексы**: по статусу, по времени старта, по году  
**Уникальность**: пара (год, сезон)

---

### CompetitionParticipant (Участник в заезде)
```go
type CompetitionParticipant struct {
    ID                string    // UUID
    CompetitionID     string    // Ссылка на Competition
    CompetitorModelID string    // Ссылка на CompetitorModel
    GridPosition      *int      // Стартовая позиция
    IsFinished        bool      // Финишировал ли
    Disqualified      bool      // Дисквалифицирован
    DNFReason         string    // Причина схода
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

**Индексы**: по заезду, по участнику

---

### CompetitionLaps (Результаты участника)
```go
type CompetitionLaps struct {
    ID                       string     // UUID
    CompetitionParticipantID string     // Ссылка на CompetitionParticipant
    TimeStart                time.Time  // Время начала
    TimeFinish               *time.Time // Время завершения
    NumberOfLaps             int        // Количество кругов
    BestLapTimeMs            int        // Лучшее время круга (мс)
    BestLapNumber            int        // Номер лучшего круга
    LastLapTimeMs            int        // Время последнего круга
    LastPassTime             *time.Time // Время последнего прохождения
    TotalCompetitionTimeMs   int        // Общее время
    CreatedAt                time.Time
    UpdatedAt                time.Time
}
```

---

### LapHistory (Детальная история кругов)
```go
type LapHistory struct {
    ID                       string    // UUID
    CompetitionParticipantID string    // Ссылка на CompetitionParticipant
    LapNumber                int       // Номер круга
    LapTimeMs                int       // Время круга (мс)
    StartTime                time.Time // Начало круга
    EndTime                  time.Time // Конец круга
    IsValid                  bool      // Валиден ли круг
    InvalidationReason       string    // Причина невалидности
    CreatedAt                time.Time
}
```

**Индексы**: по участнику, по номеру круга

---

### RawScan (Сырые RFID сканирования)
```go
type RawScan struct {
    ID                      string     // UUID
    TagValue                string     // Значение RFID метки
    Timestamp               time.Time  // Время сканирования
    ReaderType              string     // Тип считывателя
    COMPort                 string     // COM-порт
    SignalStrength          *int       // Уровень сигнала (опционально)
    IsProcessed             bool       // Обработано ли
    LinkedCompetitorModelID *string    // Связь с CompetitorModel (опционально)
    CreatedAt               time.Time
}
```

**Индексы**: по времени, по обработанному статусу, по linked_competitor_model_id

---

## Справочники

### RCModelBrand
```go
type RCModelBrand struct {
    ID        string    // UUID
    Name      string    // Название бренда (уникальное)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Примеры: XRAY, Mugen, Serpent, Associated, Yokomo

---

### RCModelScale
```go
type RCModelScale struct {
    ID        string    // UUID
    Name      string    // Масштаб (уникальный)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Примеры: 1:10, 1:8, 1:12, 1:5

---

### RCModelType
```go
type RCModelType struct {
    ID        string    // UUID
    Name      string    // Тип (уникальный)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Примеры: Touring Car, Buggy, Truggy, Monster Truck, Onroad

---

### CompetitionTrack
```go
type CompetitionTrack struct {
    ID        string    // UUID
    Name      string    // Название трассы (уникальное)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

Примеры: Main Track, Practice Track, Indoor A

---

### CompetitionYear
```go
type CompetitionYear struct {
    ID        string    // UUID
    Year      int       // Год (уникальный)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

### CompetitionSeason
```go
type CompetitionSeason struct {
    ID      string    // UUID
    Season  string    // Сезон (уникальный): Spring, Summer, Fall, Winter
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

### Settings (Системные настройки)
```go
type SystemSetting struct {
    Key         string    // Ключ настройки (уникальный)
    Value       string    // Значение
    ValueType   string    // Тип значения (string, int, bool)
    Description string    // Описание настройки
    UpdatedAt   time.Time
}
```

Используется для хранения настроек приложения (локаль, последний COM-порт и т.д.)

---

### AuditLog (Аудит действий)
```go
type AuditLog struct {
    ID         string    // UUID
    Timestamp  time.Time // Время события
    ActionType string    // Тип действия (CREATE, UPDATE, DELETE)
    EntityType string    // Тип сущности
    EntityID   string    // ID сущности
    UserName   string    // Имя пользователя
    IPAddress  string    // IP адрес
    Details    string    // Детали в JSON
    CreatedAt  time.Time
}
```

**Индексы**: по timestamp, по entity_type, по action_type

---

## Связи между таблицами

```
Competitor (1) ──< CompetitorModel >── (1) RCModel
                          │
                          │ (через transponder)
                          ▼
              CompetitionParticipant >── (1) Competition
                          │
                          ▼
                    CompetitionLaps
                          │
                          ▼
                      LapHistory

RawScan (независимая таблица сырых данных, может быть связана с CompetitorModel)
SystemSetting (независимая таблица настроек)
AuditLog (независимая таблица аудита)
```

---

## Бизнес-правила

### Competitor
- `CompetitorNumber` должен быть уникальным
- Рейтинг: целое число (обычно 0-3000)

### CompetitorModel
- Один транспондер может быть привязан только к одной активной связи
- При смене владельца транспондера старая связь деактивируется

### Competition
- Статусы: `scheduled` → `in_progress` → `finished` или `cancelled`
- Нельзя изменить заезд после начала (`in_progress`)
- Лимит времени ИЛИ лимит кругов (может быть оба)

### LapHistory
- Круг может быть невалидным (invalidation_reason заполняется)
- Причины невалидности: "false positive", "missed start", "cut track"

---

## Примеры запросов

### Получить все модели пилота
```sql
SELECT rm.* FROM rc_models rm
JOIN competitor_models cm ON rm.id = cm.rc_model_id
WHERE cm.competitor_id = ? AND cm.is_active = 1
```

### Получить участников заезда с результатами
```sql
SELECT c.full_name, rm.model_name, cl.number_of_laps, cl.best_lap_time_ms
FROM competition_participants cp
JOIN competitor_models cm ON cp.competitor_model_id = cm.id
JOIN competitors c ON cm.competitor_id = c.id
JOIN rc_models rm ON cm.rc_model_id = rm.id
LEFT JOIN competition_laps cl ON cp.id = cl.competition_participant_id
WHERE cp.competition_id = ?
ORDER BY cl.number_of_laps DESC, cl.best_lap_time_ms ASC
```

### Найти активный транспондер по номеру
```sql
SELECT cm.*, c.full_name, rm.model_name
FROM competitor_models cm
JOIN competitors c ON cm.competitor_id = c.id
JOIN rc_models rm ON cm.rc_model_id = rm.id
WHERE cm.transponder_number = ? AND cm.is_active = 1
```
