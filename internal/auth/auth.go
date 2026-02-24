package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication logic
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// IsAuthenticated checks if the application is authenticated for admin functions
func (as *AuthService) IsAuthenticated() bool {
	// For now, we'll store the authentication state in memory
	// In a real implementation, you might want to use sessions or tokens
	return isAuthenticated
}

// Authenticate validates user credentials
func (as *AuthService) Authenticate(username, password string) (bool, error) {
	var storedHash string
	err := as.db.QueryRow("SELECT value FROM system_settings WHERE key = 'auth_password_hash'").Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("authentication system not properly initialized")
		}
		return false, fmt.Errorf("failed to retrieve password hash: %w", err)
	}

	storedUsername := ""
	err = as.db.QueryRow("SELECT value FROM system_settings WHERE key = 'auth_user'").Scan(&storedUsername)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve username: %w", err)
	}

	if username != storedUsername {
		as.logAuthAttempt(username, false, "incorrect username")
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		as.logAuthAttempt(username, false, "incorrect password")
		return false, nil
	}

	as.logAuthAttempt(username, true, "")
	
	// Set the authenticated flag
	isAuthenticated = true
	return true, nil
}

// ChangeCredentials updates the username and password in the database
func (as *AuthService) ChangeCredentials(currentUsername, currentPassword, newUsername, newPassword string) error {
	// First authenticate with current credentials
	authed, err := as.Authenticate(currentUsername, currentPassword)
	if err != nil {
		return fmt.Errorf("error during authentication: %w", err)
	}
	if !authed {
		return fmt.Errorf("current credentials are invalid")
	}

	// Hash the new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update the database
	tx, err := as.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update password
	_, err = tx.Exec("UPDATE system_settings SET value = ?, updated_at = datetime('now') WHERE key = 'auth_password_hash'", string(newHash))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Update username if it's different
	if newUsername != currentUsername {
		_, err = tx.Exec("UPDATE system_settings SET value = ?, updated_at = datetime('now') WHERE key = 'auth_user'", newUsername)
		if err != nil {
			return fmt.Errorf("failed to update username: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Credentials updated successfully for user: %s", newUsername)
	as.logAuditEvent("CHANGE_CREDENTIALS", "SYSTEM", "", newUsername, "User changed their credentials")
	
	return nil
}

// logAuthAttempt logs authentication attempts to the audit log
func (as *AuthService) logAuthAttempt(username string, success bool, reason string) {
	action := "AUTH_LOGIN_SUCCESS"
	if !success {
		action = "AUTH_LOGIN_FAILED"
	}

	details := fmt.Sprintf("User: %s", username)
	if reason != "" {
		details += fmt.Sprintf(", Reason: %s", reason)
	}

	_, err := as.db.Exec(`
		INSERT INTO audit_log (timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at)
		VALUES (datetime('now'), ?, 'AUTH', '', ?, '', ?, datetime('now'))
	`, action, username, details)

	if err != nil {
		log.Printf("Failed to log authentication attempt: %v", err)
	}
}

// logAuditEvent logs an audit event
func (as *AuthService) logAuditEvent(actionType, entityType, entityID, userName, details string) {
	_, err := as.db.Exec(`
		INSERT INTO audit_log (timestamp, action_type, entity_type, entity_id, user_name, ip_address, details, created_at)
		VALUES (datetime('now'), ?, ?, ?, ?, '', ?, datetime('now'))
	`, actionType, entityType, entityID, userName, details)

	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
	}
}

// Logout resets the authentication state
func (as *AuthService) Logout() {
	isAuthenticated = false
	log.Println("User logged out")
}

// Global variable to track authentication status
// In a production application, you would want a more robust session management system
var isAuthenticated bool = false