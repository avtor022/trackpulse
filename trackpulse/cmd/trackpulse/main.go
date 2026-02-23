package main

import (
	"log"

	"trackpulse/internal/database"
	"trackpulse/internal/race"
	"trackpulse/internal/serial"
	"trackpulse/internal/websocket"
	"trackpulse/ui/desktop"
)

func main() {
	// Initialize database
	db, err := database.NewDB("trackpulse.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize race controller
	raceController := race.NewController(db)

	// Initialize serial communication
	serialHandler := serial.NewHandler(raceController)
	go serialHandler.StartListening()

	// Initialize WebSocket server
	wsServer := websocket.NewServer()
	go wsServer.Start()

	// Start desktop UI
	desktop.StartUI(db, raceController, wsServer)
}