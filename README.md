# trackpulse
RFID trackpulse app

## Запуск проекта

> **ВАЖНО:** Все команды выполняются ТОЛЬКО ИЗ КОРНЯ ПРОЕКТА (trackpulse/)!

### Подготовка

```bash
# Очистить кэш
go clean --cache

# Проверка модулей
go mod tidy

# Запуск без сборки
go run ./cmd/main.go
```

### Сборка

**Windows:**
```bash
go build -ldflags="-s -w -H windowsgui" -o trackpulse.exe cmd\main.go
```

**Linux:**
```bash
go build -ldflags="-s -w" -o trackpulse cmd/main.go
chmod +x trackpulse
```

**MacOS:**
```bash
go build -ldflags="-s -w" -o trackpulse cmd/main.go
```

### Запуск собранного приложения

```bash
./trackpulse      # Linux/Mac
trackpulse.exe    # Windows
```

## Мониторинг заезда

Создайте соревнование, привяжите к нему транспондеры участников и подключите считыватель. На вкладке «Мониторинг» выберите соревнование и нажмите «Старт». Таблица показывает место, номер и имя спортсмена, число завершённых кругов и время быстрейшего круга в формате мм:сс:мс.

Первая отметка транспондера начинает первый круг; следующая допустимая отметка завершает его. Минимальное время круга задаётся в форме соревнования в секундах, по умолчанию — 5. Отметка раньше этого времени игнорируется и не меняет начало текущего круга. Места определяются по количеству кругов, затем по лучшему кругу; при полном равенстве — по номеру спортсмена. Результаты сохраняются в базе данных.

## Тестирование

### Запуск тестов для сервиса участников

```bash
go test ./tests/service/ -v -run TestCompetitorService
```

### Запуск всех тестов в проекте

```bash
go test ./... -v
```

### Отчёт о покрытии кода

```bash
go test ./tests/service/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Запуск конкретного теста

```bash
go test ./tests/service/ -v -run TestCompetitorService_CreateCompetitor_Success
```
