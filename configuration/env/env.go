package env

import (
	"net/url"
	"os"

	"github.com/PaluMacil/dwn/configuration"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Config() configuration.Configuration {
	ws := configuration.WebServerConfiguration{
		Protocol:    os.Getenv("DWN_PROTOCOL"),
		Host:        os.Getenv("DWN_HOST"),
		Port:        os.Getenv("DWN_PORT"),
		BaseURL:     os.Getenv("DWN_BASE_URL"),
		UIProxyPort: os.Getenv("DWN_UI_PROXY_PORT"),
		ContentRoot: os.Getenv("DWN_CONTENT_ROOT"),
	}
	home, err := url.Parse(ws.HomePage())
	if err != nil {
		panic("Cannot parse home URL: " + err.Error())
	}
	googleCallbackURL, err := url.Parse("oauth/google/callback")
	if err != nil {
		panic("Cannot parse google callback URL: " + err.Error())
	}
	googleRedirect := home.ResolveReference(googleCallbackURL)
	return configuration.Configuration{
		WebServer: ws,
		Setup: configuration.SetupConfiguration{
			InitialAdmin: os.Getenv("DWN_INITIAL_ADMIN"),
		},
		Database: configuration.DatabaseConfiguration{
			DataDir: os.Getenv("DWN_DATA_DIR"),
		},
		Auth: configuration.AuthConfiguration{
			Google: &oauth2.Config{
				RedirectURL:  googleRedirect.String(),
				ClientID:     os.Getenv("DWN_OAUTH_GOOGLE_KEY"),
				ClientSecret: os.Getenv("DWN_OAUTH_GOOGLE_SECRET"),
				Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint: google.Endpoint,
			},
		},
	}
}
