package main

import (
	"log"
	"trackpulse/internal/config"
	"trackpulse/internal/repository"
	"trackpulse/internal/services"
	"trackpulse/internal/ui"
)

func main() {
	// 1. Инициализация конфигурации
	cfg := config.Load()

	// 2. Инициализация базы данных
	db, err := repository.InitDatabase(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 3. Инициализация репозиториев
	repos := repository.NewRepositories(db)

	// 4. Инициализация сервисов
	svc := services.NewServices(repos, cfg)

	// 5. Инициализация UI
	app := ui.NewApplication(svc, cfg)

	// 6. Запуск приложения
	app.Run()
}