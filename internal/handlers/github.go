package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/debridmediamanager/gatekeeper/internal/auth"
	"github.com/debridmediamanager/gatekeeper/internal/services"
	"github.com/debridmediamanager/gatekeeper/pkg/models"

	"github.com/gin-gonic/gin"
)

func GithubLoginHandler(c *gin.Context) {
	url := auth.GithubOAuthConfig.AuthCodeURL("random")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GithubCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != "random" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := auth.ExchangeGithubToken(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := auth.GithubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var user models.GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	if services.IsGitHubSponsor(user.Login) {
		c.Redirect(http.StatusTemporaryRedirect, "/discord")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/patreon")
	}
}
