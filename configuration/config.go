package configuration

import (
	"fmt"
	"golang.org/x/oauth2/google"
	"net/url"
	"os"

	"golang.org/x/oauth2"
)

type SetupConfiguration struct {
	InitialAdmin string `json:"initialAdmin"`
}

type Configuration struct {
	WebServer WebServerConfiguration `json:"webServer"`
	Setup     SetupConfiguration     `json:"setup"`
	Database  DatabaseConfiguration  `json:"database"`
	Auth      AuthConfiguration      `json:"auth"`
	Prod      bool                   `json:"prod"`
}

type WebServerConfiguration struct {
	Protocol    string `json:"protocol"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	UIProxyPort string `json:"uiProxyPort"`
	ContentRoot string `json:"contentRoot"`
}

type DatabaseConfiguration struct {
	DataDir       string `json:"dataDir"`
	EncryptionKey []byte `json:"-"`
}

func (ws WebServerConfiguration) HomePage() string {
	port := ws.UIProxyPort
	if port == "" {
		port = ws.Port
	}
	return fmt.Sprintf("%s://%s:%s",
		ws.Protocol, ws.Host, port)
}

type AuthConfiguration struct {
	Google *oauth2.Config `json:"google"`
}

func New(prod bool) (Configuration, error) {
	env := envType{prod: prod}
	const devEncKey = "3d17618d4297f83665b32e28f9b1c23d"
	var valProtocol, valHost, valPort, valUIProxyPort, valContentRoot, valInitialAdmin, valDataDir, valEncryptionKey = env.coalesce("DWN_PROTOCOL", "https", "http"),
		env.coalesce("DWN_HOST", "danwolf.net", "localhost"),
		env.coalesce("DWN_PORT", "3035", "3035"),
		env.coalesce("DWN_UI_PROXY_PORT", "443", "4200"),
		//TODO: get this last folder level generalized and usable per project
		env.coalesce("DWN_CONTENT_ROOT", "/opt/danwolf.net/dwn-ui/dist/dwn-ui", "/home/dan/repos/dwn-ui/dist/dwn-ui"),
		env.coalesce("DWN_INITIAL_ADMIN", "dcwolf@gmail.com", "dcwolf@gmail.com"),
		env.coalesce("DWN_DATA_DIR", "data", "data"),
		env.coalesce("DWN_MASTER_ENC_KEY", devEncKey, devEncKey)

	if prod && valEncryptionKey == devEncKey {
		return Configuration{}, fmt.Errorf("a encryption key must not be empty or the same as the dev key")
	}

	ws := WebServerConfiguration{
		Protocol:    valProtocol,
		Host:        valHost,
		Port:        valPort,
		UIProxyPort: valUIProxyPort,
		ContentRoot: valContentRoot,
	}
	home, err := url.Parse(ws.HomePage())
	if err != nil {
		return Configuration{}, fmt.Errorf("cannot parse home URL: %w", err)
	}
	googleCallbackURL, err := url.Parse("oauth/google/callback")
	if err != nil {
		return Configuration{}, fmt.Errorf("cannot parse google callback URL: %w", err)
	}
	googleRedirect := home.ResolveReference(googleCallbackURL)
	return Configuration{
		WebServer: ws,
		Setup: SetupConfiguration{
			InitialAdmin: valInitialAdmin,
		},
		Database: DatabaseConfiguration{
			DataDir:       valDataDir,
			EncryptionKey: []byte(valEncryptionKey),
		},
		Auth: AuthConfiguration{
			// TODO: add auth providers when and if env vars exist (or later people can use config)
			Google: &oauth2.Config{
				RedirectURL:  googleRedirect.String(),
				ClientID:     os.Getenv("DWN_OAUTH_GOOGLE_KEY"),
				ClientSecret: os.Getenv("DWN_OAUTH_GOOGLE_SECRET"),
				Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
					"https://www.googleapis.com/auth/userinfo.email"},
				Endpoint: google.Endpoint,
			},
		},
		Prod: prod,
	}, nil
}

type envType struct {
	prod bool
}

func (e envType) coalesce(envKey, prodDefault, devDefault string) string {
	value := os.Getenv(envKey)
	if value == "" {
		if e.prod {
			return prodDefault
		} else {
			return devDefault
		}
	}
	return value
}
