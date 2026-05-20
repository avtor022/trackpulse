# TrackPulse - UI Компоненты (Fyne)

## Обзор UI слоя

UI компоненты находятся в `internal/ui/` и используют фреймворк **Fyne v2.7.3**.

**Направление зависимостей**: UI → Service → Repository

---

## App (app.go)

### Назначение
Главное окно приложения, навигация по вкладкам, инициализация.

### Структура
```go
type App struct {
    app         fyne.App
    window      fyne.Window
    competitorService  *service.CompetitorService
    modelService       *service.RCModelService
    settingsService    *service.SettingsService
    competitorModelService *service.CompetitorModelService
    competitionService *service.CompetitionService
    participantService *service.CompetitionParticipantService
    
    // Панели
    competitorPanel  *CompetitorPanel
    rcModelPanel     *RCModelPanel
    competitorModelPanel *CompetitorModelPanel
    competitionPanel *CompetitionPanel
    monitoringPanel  *MonitoringPanel
    settingsPanel    *SettingsPanel
    logsPanel        *LogsPanel
    
    locale string
}
```

### Методы
```go
NewApp(...) *App
Run()                    // Запуск приложения
buildMainMenu() *fyne.MainMenu
switchTab(tabID string)  // Переключение вкладок
```

### Вкладки приложения
1. **Competitors** — управление пилотами
2. **RC Models** — каталог моделей
3. **Transponders** — привязка транспондеров
4. **Competitions** — создание заездов
5. **Monitoring** — экран гонки (real-time)
6. **Logs** — просмотр логов
7. **Settings** — настройки

---

## CompetitorPanel (competitor_panel.go)

### Назначение
Управление участниками: список, создание, редактирование, удаление.

### Элементы UI
- `table` — таблица участников (номер, имя, рейтинг, город)
- `searchEntry` — поиск по имени
- `addButton` — добавить нового
- `editButton` — редактировать выбранного
- `deleteButton` — удалить выбранного
- `formPopup` — форма создания/редактирования

### Форма участника
Поля:
- Competitor Number (автогенерация)
- Full Name (обязательно)
- Birthday (date picker)
- Country (text)
- City (text)
- Rating (number)

### Методы
```go
NewCompetitorPanel(service) *CompetitorPanel
refreshTable()           // Обновить таблицу
showAddForm()            // Показать форму добавления
showEditForm(id)         // Показать форму редактирования
deleteSelected()         // Удалить выбранного
onSearchChanged(query)   // Обработка поиска
```

---

## RCModelPanel (rc_model_panel.go)

### Назначение
Каталог RC-моделей: список, фильтры, CRUD.

### Элементы UI
- `table` — таблица моделей (бренд, название, масштаб, тип)
- `filterContainer` — панель фильтров
- `brandFilter` — выбор бренда (select)
- `scaleFilter` — выбор масштаба (select)
- `typeFilter` — выбор типа (select)
- `addButton`, `editButton`, `deleteButton`
- `formPopup` — форма модели

### Форма модели
Поля:
- Brand (autocomplete из справочника)
- Model Name (обязательно)
- Scale (select: 1:10, 1:8, 1:12...)
- Model Type (select: Touring, Buggy...)
- Motor Type (text)
- Drive Type (select: 2WD, 4WD)

### Методы
```go
NewRCModelPanel(service) *RCModelPanel
refreshTable()
applyFilters()
showAddForm()
showEditForm(id)
getAvailableBrands() []string
getAvailableScales() []string
getAvailableTypes() []string
```

---

## CompetitorModelPanel (competitor_model_panel.go)

### Назначение
Привязка транспондеров к участникам и моделям.

### Элементы UI
- `table` — список связей (пилот, модель, транспондер, статус)
- `addButton`, `editButton`, `deactivateButton`
- `formPopup` — форма привязки

### Форма привязки
Поля:
- Competitor (autocomplete с поиском)
- RC Model (autocomplete с фильтрами)
- Transponder Number (обязательно, уникально)
- Transponder Type (select: RFID, AMB)
- Is Active (checkbox)

### Методы
```go
NewCompetitorModelPanel(service, compService, modelService) *CompetitorModelPanel
refreshTable()
showAddForm()
searchCompetitor(query) []Competitor
searchModel(query) []RCModel
validateTransponder(number) bool
deactivateSelected()
```

### Автодополнение
Использует `AutocompleteEntry` для:
- Поиска пилота по имени
- Поиска модели по бренду+названию

---

## CompetitionPanel (competition_panel.go)

### Назначение
Создание и управление заездами.

### Элементы UI
- `table` — список заездов (название, тип, статус, время)
- `filterContainer` — фильтры (год, сезон, статус)
- `addButton`, `editButton`, `deleteButton`
- `startButton`, `finishButton` — управление заездом
- `participantsButton` — управление участниками
- `formPopup` — форма заезда
- `participantPopup` — выбор участников

### Форма заезда
Поля:
- Competition Title (обязательно)
- Competition Type (select: qualifying, final, practice)
- Model Type (select)
- Model Scale (select)
- Track Name (select)
- Lap Count Target (number, опционально)
- Time Limit Minutes (number, опционально)
- Competition Year (select)
- Season (select: Spring, Summer, Fall, Winter)
- Status (read-only или select)

### Управление участниками
Отдельное popup окно:
- Список доступных CompetitorModel
- Выбранные участники
- Grid Position (drag&drop или number input)
- Add/Remove кнопки

### Методы
```go
NewCompetitionPanel(service, participantService) *CompetitionPanel
refreshTable()
showAddForm()
showParticipantManager(competitionID)
startSelectedCompetition()
finishSelectedCompetition()
validateCompetition(dto) error
```

### Статусы и действия
| Статус | Доступные действия |
|--------|-------------------|
| scheduled | Edit, Delete, Start, Manage Participants |
| in_progress | Finish, Cancel, View Only |
| finished | View Only |
| cancelled | View Only |

---

## MonitoringPanel (monitoring_panel.go)

### Назначение
Экран гонки в реальном времени: таймер, круги, результаты, регистрация участников.

### Элементы UI
- `timerLabel` — текущее время заезда (MM:SS)
- `statusLabel` — статус (Running, Finished)
- `competitionButton` — выбор соревнования
- `startButton`, `stopButton` — управление гонкой
- `participantsTable` — таблица регистрации участников с данными о кругах
- `competitionFilter` — фильтр соревнований

### Таблица регистрации участников
Отображается сразу после выбора соревнования. Обновляется в реальном времени.

Колонки:
- **Транспондер** — работоспособность (✓/✗): false по умолчанию, true после первого проезда через антенну
- **Номер** — номер участника
- **ФИО** — полное имя участника
- **Модель** — название модели
- **Масштаб** — масштаб модели
- **Круги** — количество пройденных кругов
- **Лучший круг** — время лучшего круга (MM:SS.mmm)

### Структура MonitoringPanel
```go
type MonitoringPanel struct {
    content                *fyne.Container
    mainWindow             fyne.Window
    competitionService     *service.CompetitionService
    participantService     *service.CompetitionParticipantService
    selectedCompetition    string
    selectedCompetitionID  string
    statusLabel            *widget.Label
    competitionButton      *widget.Button
    startButton            *widget.Button
    stopButton             *widget.Button
    timerLabel             *widget.Label
    timer                  *Timer
    competitionFilter      *CompetitionFilter
    participantsTable      *widget.Table
    participantsContainer  *fyne.Container
}
```

### Методы
```go
NewMonitoringPanel(competitionService, participantService, mainWindow) *MonitoringPanel
createContent() *fyne.Container
createParticipantsTable() *fyne.Container
updateParticipantsTable()                        // Обновление таблицы участников
onCompetitionSelected(selected string)           // Обработка выбора соревнования
startMonitoring()                                // Старт гонки
stopMonitoring()                                 // Остановка гонки
refreshCompetitions()                            // Обновление списка соревнований
UpdateData()                                     // Обновление данных панели
Refresh()                                        // Перерисовка панели
```

### Интеграция с LapService
- При первом проезде участника через антенну LapService вызывает `markTransponderWorked(participantID)`
- Флаг `TransponderWorked` сохраняется в БД и отображается в таблице как "✓"
- Данные для таблицы загружаются через `participantService.GetParticipantRegistrationData(competitionID)`
- Таблица обновляется автоматически при выборе соревнования и периодически во время гонки

### Таймер
Компонент `Timer`:
- Обратный отсчёт (если установлен лимит)
- Прямой отсчёт от старта
- Пауза/продолжение
- Сигнал окончания

---

## SettingsPanel (settings_panel.go)

### Назначение
Настройки приложения.

### Элементы UI
- `localeSelect` — выбор языка (en, ru)
- `comPortSelect` — выбор COM-порта
- `portScanButton` — сканирование портов
- `debounceInput` — debounce время (ms)
- `readerTypeSelect` — тип считывателя
- `saveButton`, `resetButton`

### Методы
```go
NewSettingsPanel(service) *SettingsPanel
scanPorts() []string
saveSettings()
loadSettings()
validateSettings() error
```

### Настройки
- **Locale**: применяется немедленно через `locale.SetLocale()`
- **COM Port**: сохраняется для следующего запуска
- **Debounce**: минимальное время между сканированиями

---

## LogsPanel (logs_panel.go)

### Назначение
Просмотр логов приложения.

### Элементы UI
- `logText` — многострочное текстовое поле
- `refreshButton` — обновить
- `clearButton` — очистить
- `levelFilter` — фильтр по уровню (INFO, ERROR, DEBUG)
- `autoScrollCheck` — автопрокрутка

### Методы
```go
NewLogsPanel() *LogsPanel
loadLogs()
parseLogFile(path) []LogEntry
filterLogs(level string) []LogEntry
```

### Формат лога
```
INFO: 2025/01/15 10:30:45 app.go:42: TrackPulse starting...
ERROR: 2025/01/15 10:30:46 database.go:23: Failed to open database
DEBUG: 2025/01/15 10:30:47 lap_service.go:156: Scan processed: TAG123
```

---

## Вспомогательные компоненты

### AutocompleteEntry (autocomplete_entry.go)
Расширенное поле ввода с автодополнением.

```go
type AutocompleteEntry struct {
    widget.Entry
    suggestions []string
    popup       *widget.PopUp
    onSelect    func(string)
}

func NewAutocompleteEntry(suggestions []string) *AutocompleteEntry
func (a *AutocompleteEntry) SetSuggestions([]string)
func (a *AutocompleteEntry) showPopup()
```

**Использование**:
- Поиск пилота
- Поиск модели
- Выбор бренда/трассы

---

### CompetitionFilter (competition_filter.go)
Панель фильтров для заездов.

Поля:
- Year (select)
- Season (select)
- Status (select)
- Model Type (select)
- Search (text)

Методы:
```go
GetFilterCriteria() FilterCriteria
Reset()
OnFilterChanged(callback)
```

---

### PortScanner (port_scanner.go)
Сканирование доступных COM-портов.

```go
ScanPorts() []PortInfo
type PortInfo struct {
    Name        string
    ProductName string
    IsAvailable bool
}
```

Использует библиотеку `go.bug.st/serial`.

---

### Timer (timer.go)
Таймер для заездов.

```go
type Timer struct {
    startTime    time.Time
    elapsedTime  time.Duration
    isRunning    bool
    limitMinutes int
    tickChan     chan time.Duration
}

func (t *Timer) Start()
func (t *Timer) Stop()
func (t *Timer) Reset()
func (t *Timer) GetElapsed() time.Duration
func (t *Timer) IsFinished() bool
```

---

### ReferencePopup (reference_popup.go)
Всплывающие окна справочников.

Используется для:
- Выбора бренда
- Выбора масштаба
- Выбора типа модели
- Выбора трассы

```go
ShowReferencePopup(title string, items []string, onSelect func(string))
```

---

## Паттерны Fyne

### Создание виджетов
```go
// Label
label := widget.NewLabel(locale.Get("key"))

// Entry
entry := widget.NewEntry()
entry.SetPlaceHolder(locale.Get("placeholder"))

// Button
button := widget.NewButton(locale.Get("action"), func() {
    // handler
})

// Table
table := widget.NewTable(
    func() (int, int) { return rows, cols },
    func() fyne.CanvasObject { return createCell() },
    func(id widget.TableCellID, obj fyne.CanvasObject) { updateCell(obj, id) },
)
```

### Layouts
```go
// VBox (вертикальный)
container.NewVBox(children...)

// HBox (горизонтальный)
container.NewHBox(children...)

// Border (границы)
container.NewBorder(top, bottom, left, right, content)

// Grid (сетка)
container.NewGridWithColumns(2, children...)

// Scroll (прокрутка)
container.NewScroll(content)
```

### Popups
```go
// Modal popup
popup := widget.NewModalPopUp(content, canvas)
popup.Show()

// centered popup
popup := widget.NewPopUp(content, canvas)
popup.ShowAtPosition(fyne.NewPos(x, y))
```

### Data binding (опционально)
```go
// Text binding
text := binding.NewString()
entry.Bind(text)

// Listen for changes
text.AddListener(binding.DataListenerFunc(func() {
    val, _ := text.Get()
    // handle change
}))
```

---

## Локализация в UI

Все тексты используются через `locale.T()`:

```go
// Labels
widget.NewLabel(locale.T("competitor.list.title"))

// Buttons
widget.NewButton(locale.T("action.add"), handler)

// Placeholders
entry.SetPlaceHolder(locale.T("field.name.placeholder"))

// Table headers
tableHeaders := []string{
    locale.T("competitor.number"),
    locale.T("competitor.name"),
    locale.T("competitor.rating"),
}
```

Ключи локализации хранятся в `internal/locale/translations/{en,ru}.json`.

---

## Обновление UI

### В главном потоке
Fyne требует обновления UI в главном потоке:

```go
// Неправильно (из goroutine)
table.Refresh()  // Может вызвать панику

// Правильно
fyne.Do(func() {
    table.Refresh()
})

// Или в main goroutine
go func() {
    // background work
    result := processData()
    
    fyne.Do(func() {
        label.SetText(result)
    })
}()
```

### Периодическое обновление
```go
ticker := time.NewTicker(1 * time.Second)
go func() {
    for range ticker.C {
        fyne.Do(func() {
            updateUI()
        })
    }
}()
```

---

## Обработка ошибок в UI

```go
func (p *CompetitorPanel) saveCompetitor(dto *CompetitorDTO) {
    competitor, err := p.service.CreateCompetitor(dto)
    if err != nil {
        dialog.ShowError(err, p.window)
        return
    }
    
    dialog.ShowInformation(
        locale.T("success"),
        locale.T("competitor.created"),
        p.window,
    )
    p.refreshTable()
}
```

---

## Советы для AI-агентов

1. **Всегда используйте `locale.T()`** для текстов
2. **Обновляйте UI через `fyne.Do()`** из горутин
3. **Кэшируйте списки** (бренды, типы) для autocomplete
4. **Обрабатывайте ошибки** через диалоги
5. **Следуйте существующим паттернам** в других панелях
6. **Тестируйте на разных размерах окна** (Fyne адаптивный)
7. **Используйте стандартные виджеты Fyne** вместо кастомных где возможно
