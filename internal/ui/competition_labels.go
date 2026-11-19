package ui

import "trackpulse/internal/locale"

func competitionTypeLabel(value string) string {
	switch value {
	case "qualifying", "final", "practice", "heat":
		return locale.T("competition.type." + value)
	default:
		return value
	}
}

func competitionStatusLabel(value string) string {
	switch value {
	case "in_progress":
		return locale.T("competition.status.running")
	case "scheduled", "finished", "cancelled":
		return locale.T("competition.status." + value)
	default:
		return value
	}
}
