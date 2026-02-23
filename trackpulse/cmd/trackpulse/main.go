package main

import (
	"log"

	"trackpulse/internal/database"
	"trackpulse/internal/ui/desktop"
)

func main() {
	// Initialize database
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Start the desktop UI
	app := desktop.NewApp(db)
	if err := app.Run(); err != nil {
		log.Fatal("Failed to run application:", err)
	}
}