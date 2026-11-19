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

## Язык интерфейса

На вкладке «Настройки» можно выбрать русский или английский язык. Выбор сохраняется в базе данных и применяется при следующем запуске. Переводы находятся в `internal/locale/translations/` и включаются в исполняемый файл при сборке. При изменении строк обновляйте оба файла переводов.
