package repository

import "database/sql"

// Repositories holds all repository instances
type Repositories struct {
	Racers           RacersRepository
	RCModels         RCModelsRepository
	RacerModels      RacerModelsRepository
	Races            RacesRepository
	RaceParticipants RaceParticipantsRepository
	RaceLaps         RaceLapsRepository
	LapHistory       LapHistoryRepository
	RawScans         RawScansRepository
	SystemSettings   SystemSettingsRepository
	AuditLog         AuditLogRepository
}

// NewRepositories creates and returns all repositories
func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Racers:           NewRacersRepository(db),
		RCModels:         NewRCModelsRepository(db),
		RacerModels:      NewRacerModelsRepository(db),
		Races:            NewRacesRepository(db),
		RaceParticipants: NewRaceParticipantsRepository(db),
		RaceLaps:         NewRaceLapsRepository(db),
		LapHistory:       NewLapHistoryRepository(db),
		RawScans:         NewRawScansRepository(db),
		SystemSettings:   NewSystemSettingsRepository(db),
		AuditLog:         NewAuditLogRepository(db),
	}
}