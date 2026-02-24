package main

import (
	"log"

	"trackpulse/internal/auth"
	"trackpulse/internal/database"
)

func main() {
	// Initialize the database
	db, err := database.InitializeDB("./track_pulse.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize the authentication service
	authService := auth.NewAuthService(db)

	// For now, just test that the DB and auth system work
	log.Println("TrackPulse application started successfully!")
	log.Println("Database initialized and ready.")
	log.Println("Authentication system ready.")

	// Test authentication with default credentials
	authenticated, err := authService.Authenticate("admin", "admin")
	if err != nil {
		log.Printf("Error authenticating: %v", err)
	} else if authenticated {
		log.Println("Successfully authenticated with default admin credentials")
	} else {
		log.Println("Authentication failed with default credentials")
	}

	// Demo functionality
	log.Println("Demo completed successfully!")
}
}