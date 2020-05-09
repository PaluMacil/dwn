package repo

import (
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/configuration"
	"golang.org/x/oauth2/google"
	"log"
	"net/url"
	"os"
	"sync"

	"golang.org/x/oauth2"
)

type ConfigurationRepo struct {
	store   database.Storer
	db      *database.Database
	current *configuration.Configuration
	lock    sync.RWMutex
}

func NewConfigurationRepo(prod bool) (*ConfigurationRepo, error) {
	config, err := initialConfiguration(prod)
	if err != nil {
		return nil, fmt.Errorf("creating initial %t configuration: %w", prod, err)
	}
	return &ConfigurationRepo{nil, nil, config, sync.RWMutex{}}, nil
}

func initialConfiguration(prod bool) (*configuration.Configuration, error) {
	mode := configuration.Mode{Prod: prod}
	const devEncKey = "3d17618d4297f83665b32e28f9b1c23d"
	var valProtocol, valHost, valPort, valUIProxyPort, valContentRoot, valInitialAdmin, valInitialPassword, valDataDir, valEncryptionKey = mode.Coalesce("DWN_PROTOCOL", "https", "http"),
		mode.Coalesce("DWN_HOST", "danwolf.net", "localhost"),
		mode.Coalesce("DWN_PORT", "3035", "3035"),
		mode.Coalesce("DWN_UI_PROXY_PORT", "443", "4200"),
		mode.Coalesce("DWN_CONTENT_ROOT", "/opt/danwolf.net/ui/", "/home/dan/repos/dwn-ui/dist/"),
		mode.Coalesce("DWN_INITIAL_ADMIN", "", "dcwolf@gmail.com"),
		mode.Coalesce("DWN_INITIAL_PASSWORD", "", ""),
		mode.Coalesce("DWN_DATA_DIR", "data", "data"),
		mode.Coalesce("DWN_MASTER_ENC_KEY", devEncKey, devEncKey)

	if prod && valEncryptionKey == devEncKey {
		return nil, fmt.Errorf("a encryption key must not be empty or the same as the dev key")
	}

	ws := configuration.WebServerConfiguration{
		Protocol:    valProtocol,
		Host:        valHost,
		Port:        valPort,
		UIProxyPort: valUIProxyPort,
		ContentRoot: valContentRoot,
	}
	home, err := url.Parse(ws.HomePage())
	if err != nil {
		return nil, fmt.Errorf("cannot parse home URL: %w", err)
	}
	googleCallbackURL, err := url.Parse("oauth/google/callback")
	if err != nil {
		return nil, fmt.Errorf("cannot parse google callback URL: %w", err)
	}
	googleRedirect := home.ResolveReference(googleCallbackURL)
	config := configuration.Configuration{
		WebServer: ws,
		Setup: configuration.SetupConfiguration{
			InitialAdmin:    valInitialAdmin,
			InitialPassword: valInitialPassword,
		},
		Database: configuration.DatabaseConfiguration{
			DataDir:       valDataDir,
			EncryptionKey: []byte(valEncryptionKey),
		},
		Prod: prod,
	}

	// Authentication Providers
	var valOAuthGoogleKey, valOAuthGoogleSecret = os.Getenv("DWN_OAUTH_GOOGLE_KEY"), os.Getenv("DWN_OAUTH_GOOGLE_SECRET")
	if valOAuthGoogleKey != "" && valOAuthGoogleSecret != "" {
		log.Println("setting Google auth provider from environment")
		config.FS.Auth.Google = &oauth2.Config{
			RedirectURL:  googleRedirect.String(),
			ClientID:     valOAuthGoogleKey,
			ClientSecret: valOAuthGoogleSecret,
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: google.Endpoint,
		}
	}

	// Other foreign systems
	var valSMTPSendGridKey, valSMTPSendGridSecret = os.Getenv("DWN_SMTP_SENDGRID_KEY"), os.Getenv("DWN_SMTP_SENDGRID_SECRET")
	if valSMTPSendGridKey != "" && valSMTPSendGridSecret != "" {
		log.Println("setting SendGrid SMTP foreign system from environment")
		config.FS.SMTP.SendGrid = &configuration.Credential{
			ID:     valSMTPSendGridKey,
			Secret: valSMTPSendGridSecret,
		}
	}

	return &config, nil
}

func (cr *ConfigurationRepo) ConnectDatabase(store database.Storer, db *database.Database) {
	// TODO: get the env creds here, get end user id, and update the creds from db storage
	cr.store, cr.db = store, db
}

func (cr *ConfigurationRepo) Get() configuration.Configuration {
	cr.lock.RLock()
	defer cr.lock.RUnlock()
	return *cr.current
}

func (cr *ConfigurationRepo) Set(config configuration.Configuration) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.current = &config
}
