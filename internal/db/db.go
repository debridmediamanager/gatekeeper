package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./db/sponsors.db")
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}

	// Create GitHubSponsors table with sponsorship details
	createGitHubSponsorsTable := `
	CREATE TABLE IF NOT EXISTS github_sponsors (
		username TEXT PRIMARY KEY,
		email TEXT,
		active_tier_amount INTEGER,
		lifetime_payments INTEGER,
		sponsorship_start_date DATETIME,
		sponsorship_end_date DATETIME,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create PatreonSponsors table with sponsorship details
	createPatreonSponsorsTable := `
	CREATE TABLE IF NOT EXISTS patreon_sponsors (
		username TEXT PRIMARY KEY,
		email TEXT,
		active_tier_amount INTEGER,
		lifetime_payments INTEGER,
		sponsorship_start_date DATETIME,
		sponsorship_end_date DATETIME,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create DiscordAccounts table
	createDiscordAccountsTable := `
	CREATE TABLE IF NOT EXISTS discord_accounts (
		username TEXT PRIMARY KEY,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create Sponsorships table
	createSponsorshipsTable := `
	CREATE TABLE IF NOT EXISTS sponsorships (
		id TEXT PRIMARY KEY,
		github_sponsor_username TEXT,
		patreon_sponsor_username TEXT,
		discord_account_username TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
		date_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (github_sponsor_username) REFERENCES github_sponsors(username),
		FOREIGN KEY (patreon_sponsor_username) REFERENCES patreon_sponsors(username),
		FOREIGN KEY (discord_account_username) REFERENCES discord_accounts(username),
		UNIQUE (github_sponsor_username, patreon_sponsor_username)
	);`

	// Execute table creation statements
	_, err = DB.Exec(createGitHubSponsorsTable)
	if err != nil {
		log.Fatalf("Failed to create the github_sponsors table: %v", err)
	}

	_, err = DB.Exec(createPatreonSponsorsTable)
	if err != nil {
		log.Fatalf("Failed to create the patreon_sponsors table: %v", err)
	}

	_, err = DB.Exec(createDiscordAccountsTable)
	if err != nil {
		log.Fatalf("Failed to create the discord_accounts table: %v", err)
	}

	_, err = DB.Exec(createSponsorshipsTable)
	if err != nil {
		log.Fatalf("Failed to create the sponsorships table: %v", err)
	}
}

// Structs representing the tables
type GitHubSponsor struct {
	Username             string
	Email                sql.NullString
	CurrentTierAmount    sql.NullInt64
	LifetimePayments     sql.NullInt64
	SponsorshipStartDate sql.NullTime
	SponsorshipEndDate   sql.NullTime
	DateCreated          time.Time
	DateUpdated          time.Time
}

type PatreonSponsor struct {
	Username             string
	Email                sql.NullString
	CurrentTierAmount    sql.NullInt64
	LifetimePayments     sql.NullInt64
	SponsorshipStartDate sql.NullTime
	SponsorshipEndDate   sql.NullTime
	DateCreated          time.Time
	DateUpdated          time.Time
}

type DiscordAccount struct {
	Username    string
	DateCreated time.Time
	DateUpdated time.Time
}

type Sponsorship struct {
	ID                     string
	GitHubSponsorUsername  sql.NullString
	PatreonSponsorUsername sql.NullString
	DiscordAccountUsername sql.NullString
	DateCreated            time.Time
	DateUpdated            time.Time
}

// Helper function to generate a UUID.
func generateUUID() string {
	return uuid.New().String()
}

// Functions to interact with the tables

// GitHubSponsors functions
func GetGitHubSponsorByUsername(username string) (*GitHubSponsor, error) {
	query := `
	SELECT username, email, active_tier_amount, lifetime_payments, sponsorship_start_date, sponsorship_end_date, date_created, date_updated
	FROM github_sponsors
	WHERE username = ?`

	row := DB.QueryRow(query, username)

	var sponsor GitHubSponsor
	err := row.Scan(
		&sponsor.Username,
		&sponsor.Email,
		&sponsor.CurrentTierAmount,
		&sponsor.LifetimePayments,
		&sponsor.SponsorshipStartDate,
		&sponsor.SponsorshipEndDate,
		&sponsor.DateCreated,
		&sponsor.DateUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No matching sponsor found.
		}
		return nil, err
	}

	return &sponsor, nil
}

func InsertGitHubSponsor(sponsor *GitHubSponsor) error {
	query := `
	INSERT INTO github_sponsors (username, email, active_tier_amount, lifetime_payments, sponsorship_start_date, sponsorship_end_date)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query, sponsor.Username, sponsor.Email, sponsor.CurrentTierAmount, sponsor.LifetimePayments, sponsor.SponsorshipStartDate, sponsor.SponsorshipEndDate)
	if err != nil {
		return err
	}

	return nil
}

func UpdateGitHubSponsor(sponsor *GitHubSponsor) error {
	query := `
	UPDATE github_sponsors
	SET email = ?, active_tier_amount = ?, lifetime_payments = ?, sponsorship_start_date = ?, sponsorship_end_date = ?, date_updated = CURRENT_TIMESTAMP
	WHERE username = ?`

	_, err := DB.Exec(query, sponsor.Email, sponsor.CurrentTierAmount, sponsor.LifetimePayments, sponsor.SponsorshipStartDate, sponsor.SponsorshipEndDate, sponsor.Username)
	return err
}

// PatreonSponsors functions
func GetPatreonSponsorByUsername(username string) (*PatreonSponsor, error) {
	query := `
	SELECT username, email, active_tier_amount, lifetime_payments, sponsorship_start_date, sponsorship_end_date, date_created, date_updated
	FROM patreon_sponsors
	WHERE username = ?`

	row := DB.QueryRow(query, username)

	var sponsor PatreonSponsor
	err := row.Scan(
		&sponsor.Username,
		&sponsor.Email,
		&sponsor.CurrentTierAmount,
		&sponsor.LifetimePayments,
		&sponsor.SponsorshipStartDate,
		&sponsor.SponsorshipEndDate,
		&sponsor.DateCreated,
		&sponsor.DateUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No matching sponsor found.
		}
		return nil, err
	}

	return &sponsor, nil
}

func InsertPatreonSponsor(sponsor *PatreonSponsor) error {
	query := `
	INSERT INTO patreon_sponsors (username, email, active_tier_amount, lifetime_payments, sponsorship_start_date, sponsorship_end_date)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(query, sponsor.Username, sponsor.Email, sponsor.CurrentTierAmount, sponsor.LifetimePayments, sponsor.SponsorshipStartDate, sponsor.SponsorshipEndDate)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePatreonSponsor(sponsor *PatreonSponsor) error {
	query := `
	UPDATE patreon_sponsors
	SET email = ?, active_tier_amount = ?, lifetime_payments = ?, sponsorship_start_date = ?, sponsorship_end_date = ?, date_updated = CURRENT_TIMESTAMP
	WHERE username = ?`

	_, err := DB.Exec(query, sponsor.Email, sponsor.CurrentTierAmount, sponsor.LifetimePayments, sponsor.SponsorshipStartDate, sponsor.SponsorshipEndDate, sponsor.Username)
	return err
}

// DiscordAccounts functions
func GetDiscordAccountByUsername(username string) (*DiscordAccount, error) {
	query := `
	SELECT username, date_created, date_updated
	FROM discord_accounts
	WHERE username = ?`

	row := DB.QueryRow(query, username)

	var account DiscordAccount
	err := row.Scan(&account.Username, &account.DateCreated, &account.DateUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No matching account found.
		}
		return nil, err
	}

	return &account, nil
}

func InsertDiscordAccount(username string) error {
	query := `
	INSERT INTO discord_accounts (username)
	VALUES (?)`

	_, err := DB.Exec(query, username)
	if err != nil {
		return err
	}

	return nil
}

// Sponsorships functions
func InsertSponsorship(sponsorship *Sponsorship) error {
	id := generateUUID()
	query := `
	INSERT INTO sponsorships (
		id,
		github_sponsor_username,
		patreon_sponsor_username,
		discord_account_username
	) VALUES (?, ?, ?, ?)`

	_, err := DB.Exec(query,
		id,
		sponsorship.GitHubSponsorUsername,
		sponsorship.PatreonSponsorUsername,
		sponsorship.DiscordAccountUsername,
	)
	if err != nil {
		return err
	}

	sponsorship.ID = id
	return nil
}

func UpdateSponsorship(sponsorship *Sponsorship) error {
	query := `
	UPDATE sponsorships SET
		github_sponsor_username = ?,
		patreon_sponsor_username = ?,
		discord_account_username = ?,
		date_updated = CURRENT_TIMESTAMP
	WHERE id = ?`

	_, err := DB.Exec(query,
		sponsorship.GitHubSponsorUsername,
		sponsorship.PatreonSponsorUsername,
		sponsorship.DiscordAccountUsername,
		sponsorship.ID,
	)
	return err
}

func GetSponsorshipByGitHubAndPatreonUsernames(githubUsername, patreonUsername string) (*Sponsorship, error) {
	query := `
	SELECT id, github_sponsor_username, patreon_sponsor_username, discord_account_username, date_created, date_updated
	FROM sponsorships
	WHERE github_sponsor_username = ? AND patreon_sponsor_username = ?`

	row := DB.QueryRow(query, githubUsername, patreonUsername)

	var sponsorship Sponsorship
	err := row.Scan(
		&sponsorship.ID,
		&sponsorship.GitHubSponsorUsername,
		&sponsorship.PatreonSponsorUsername,
		&sponsorship.DiscordAccountUsername,
		&sponsorship.DateCreated,
		&sponsorship.DateUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No matching sponsorship found.
		}
		return nil, err
	}

	return &sponsorship, nil
}
