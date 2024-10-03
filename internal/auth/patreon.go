package auth

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

var PatreonOAuthConfig *oauth2.Config

func InitPatreonOAuth() {
	PatreonOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("PATREON_CLIENT_ID"),
		ClientSecret: os.Getenv("PATREON_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/patreon/callback",
		Scopes:       []string{"identity", "identity.memberships"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.patreon.com/oauth2/authorize",
			TokenURL: "https://www.patreon.com/api/oauth2/token",
		},
	}
}

func ExchangePatreonToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return PatreonOAuthConfig.Exchange(ctx, code)
}
