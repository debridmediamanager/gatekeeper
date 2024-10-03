package main

import (
	"log"

	"github.com/debridmediamanager/gatekeeper/internal/auth"
	"github.com/debridmediamanager/gatekeeper/internal/db"
	"github.com/debridmediamanager/gatekeeper/internal/handlers"
	"github.com/debridmediamanager/gatekeeper/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()
	db.InitDB()

	auth.InitGithubOAuth()

	router := gin.Default()
	router.LoadHTMLGlob("web/*")
	router.Static("/static", "./static")

	// Routes
	router.GET("/", handlers.IndexHandler)
	router.GET("/begin", handlers.BeginHandler)
	router.GET("/github", handlers.GitHubHandler)
	router.GET("/patreon", handlers.PatreonHandler)
	router.GET("/discord", handlers.DiscordHandler)
	router.GET("/end", handlers.EndHandler)
	// fails
	router.GET("/already-mapped", handlers.AlreadyMappedHandler)
	router.GET("/not-enough", handlers.NotEnoughHandler)
	router.GET("/not-sponsor", handlers.NotSponsorHandler)

	// Github OAuth
	router.GET("/auth/github", handlers.GithubLoginHandler)
	router.GET("/auth/github/callback", handlers.GithubCallbackHandler)

	// Patreon OAuth
	router.GET("/auth/patreon", handlers.PatreonLoginHandler)
	router.GET("/auth/patreon/callback", handlers.PatreonCallbackHandler)

	// Discord OAuth
	router.GET("/auth/discord", handlers.DiscordLoginHandler)
	router.GET("/auth/discord/callback", handlers.DiscordCallbackHandler)

	log.Fatal(router.Run(":8080"))
}
