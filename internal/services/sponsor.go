package services

import (
	"fmt"

	"github.com/debridmediamanager/gatekeeper/internal/db"
)

func GrantAccessAndSave(githubUsername, patreonUsername string, tierAmount, lifetimePayments int) error {
	stmt, err := db.DB.Prepare(`
		INSERT INTO sponsors (github_username, patreon_username, tier_amount, lifetime_payments, date_updated)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(githubUsername, patreonUsername, tierAmount, lifetimePayments)
	if err != nil {
		return fmt.Errorf("failed to insert sponsor: %w", err)
	}

	return nil
}

func IsGitHubSponsor(username string) bool {
	return true
}
