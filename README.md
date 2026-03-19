# trackpulse
RFID trackpulse app

## Запуск проекта

# ТОЛЬКО ИЗ КОРНЯ ПРОЕКТА (trackpulse/) ! 
# ✅ очистить кэш
go clean --cache

# ✅ Проверка модулей 
go mod tidy 

# ✅ Запуск без сборки
go run ./cmd/main.go

# ✅ Сборка 
go build -ldflags="-s -w -H windowsgui" -o trackpulse.exe cmd\main.go

# ✅ Запуск 
./trackpulse # Linux/Mac 
trackpulse.exe # Windows
