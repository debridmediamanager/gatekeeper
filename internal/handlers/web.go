package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func BeginHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/github")
}

func GitHubHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "github.html", nil)
}

func PatreonHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "patreon.html", nil)
}

func DiscordHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "discord.html", nil)
}

func EndHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "end.html", nil)
}

func AlreadyMappedHandler(c *gin.Context) {
	c.HTML(http.StatusForbidden, "already_mapped.html", nil)
}

func NotEnoughHandler(c *gin.Context) {
	c.HTML(http.StatusForbidden, "not_enough.html", nil)
}

func NotSponsorHandler(c *gin.Context) {
	c.HTML(http.StatusForbidden, "not_sponsor.html", nil)
}
