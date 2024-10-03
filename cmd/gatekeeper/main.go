package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuth configurations
var (
	githubOAuthConfig  *oauth2.Config
	patreonOAuthConfig *oauth2.Config
	oauthStateString   = "random" // Replace with a random string generator
)

// Database
var db *sql.DB

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize GitHub OAuth configuration
	githubOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		Scopes:       []string{"read:user"},
		Endpoint:     github.Endpoint,
	}

	// Initialize Patreon OAuth configuration
	patreonOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("PATREON_CLIENT_ID"),
		ClientSecret: os.Getenv("PATREON_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/patreon/callback",
		Scopes:       []string{"identity", "identity.memberships"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.patreon.com/oauth2/authorize",
			TokenURL: "https://www.patreon.com/api/oauth2/token",
		},
	}

	// Initialize SQLite database
	db, err = sql.Open("sqlite3", "./db/sponsors.db")
	if err != nil {
		panic(err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS sponsors (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        github_username TEXT,
        patreon_username TEXT,
        tier_amount INTEGER,
        lifetime_payments INTEGER,
        date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
        date_updated DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err = db.Exec(createTable); err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("web/*")
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// GitHub OAuth routes
	router.GET("/login/github", githubLoginHandler)
	router.GET("/auth/github/callback", githubCallbackHandler)

	// Patreon OAuth routes
	router.GET("/login/patreon", patreonLoginHandler)
	router.GET("/auth/patreon/callback", patreonCallbackHandler)

	router.Run(":8080")
}

func githubLoginHandler(c *gin.Context) {
	url := githubOAuthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func githubCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != oauthStateString {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Check if the user is a sponsor
	if isGitHubSponsor(user.Login) {
		log.Println("User is a GitHub sponsor")
		// Proceed to Discord login (hidden at first)
		c.HTML(http.StatusOK, "discord_login.html", nil)
	} else {
		log.Println("User is not a GitHub sponsor")
		// Show Patreon login button
		c.HTML(http.StatusOK, "patreon_login.html", nil)
	}
}

func isGitHubSponsor(username string) bool {
	// Implement GitHub Sponsors API call here
	// Note: GitHub Sponsors API requires special access
	// For demonstration, let's assume the function returns false
	return true
}

func patreonLoginHandler(c *gin.Context) {
	url := patreonOAuthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func patreonCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != oauthStateString {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := patreonOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://www.patreon.com/api/oauth2/v2/identity?include=memberships")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var user struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				FullName string `json:"full_name"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Check if the user is a Patreon sponsor
	if isPatreonSponsor(user.Data.ID) {
		// Grant access and save to database
		go grantAccessAndSave(user.Data.Attributes.FullName, "", 0, 0)
		c.HTML(http.StatusOK, "success.html", nil)
	} else {
		// Show message to sponsor the project
		c.HTML(http.StatusOK, "sponsor_message.html", nil)
	}
}

func isPatreonSponsor(userID string) bool {
	// Implement Patreon API call to check sponsorship
	// For demonstration, let's assume the function returns true
	return true
}

func grantAccessAndSave(githubUsername, patreonUsername string, tierAmount, lifetimePayments int) {
	// Grant access to the private GitHub repository
	if githubUsername != "" {
		addUserToRepo(githubUsername)
	}

	// Save user data to SQLite database
	stmt, err := db.Prepare(`
        INSERT INTO sponsors (github_username, patreon_username, tier_amount, lifetime_payments, date_updated)
        VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
    `)
	if err != nil {
		fmt.Println("Database error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(githubUsername, patreonUsername, tierAmount, lifetimePayments)
	if err != nil {
		fmt.Println("Database error:", err)
	}
}

func addUserToRepo(username string) {
	accessToken := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	url := fmt.Sprintf("https://api.github.com/repos/debridmediamanager/zurg/collaborators/%s", username)

	req, _ := http.NewRequest("PUT", url, nil)
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent) {
		fmt.Println("Failed to add collaborator:", err)
	}
}
