package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/core"
	"github.com/PaluMacil/dwn/database"
)

type Module struct {
	db     *database.Database
	config configuration.Configuration
}

func New(db *database.Database, config configuration.Configuration) *Module {
	return &Module{
		db:     db,
		config: config,
	}
}

func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/") //route[2]) is google
	if len(route) != 4 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch endpoint := route[3]; endpoint {
	case "login":
		url := mod.config.Auth.Google.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	case "callback":
		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := mod.config.Auth.Google.Exchange(r.Context(), code)
		if err != nil {
			fmt.Printf("Code exchange failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer response.Body.Close()
		claims := &GoogleClaims{}
		err = json.NewDecoder(response.Body).Decode(claims)
		if err != nil {
			http.Error(w, "Google Claims: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if claims.VerifiedEmail {
			ip := r.RemoteAddr
			forwardForIP := r.Header.Get("X-Forwarded-For")
			if forwardForIP != "" {
				// if the user if coming through one or more proxies, one or more IP addresses
				// could be set in this comma separated header. The first IP is the user's
				// original IP.
				ip = strings.Split(forwardForIP, ",")[0]
			}
			session := mod.db.Sessions.GenerateFor(claims.Email, ip)
			err := mod.db.Sessions.Set(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//if user exists in database, save session, update last login
			user, err := mod.db.Users.Get(claims.Email)
			if mod.db.Util.IsKeyNotFoundErr(err) {
				displayName, err := generateDisplayName(claims.GivenName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user = claims.CreateUser(displayName)
				if claims.Email == mod.config.Setup.InitialAdmin {
					//TODO: handle err below and add other users to User group
					mod.db.UserGroups.Set(core.UserGroup{
						Email:     claims.Email,
						GroupName: core.BuiltInGroupAdmin,
					})
				}
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user.LastLogin = time.Now()
			err = mod.db.Users.Set(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl, err := template.New("login").Parse(loginCallbackPage)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w,
				loginCallbackData{
					TokenName:   "core-token",
					Token:       session.Token,
					RedirectURL: mod.config.WebServer.HomePage(), //TODO: Check to see if a different redirect is requested and if it is safe
				})
			return
		} else {
			//TODO: Send to registration page. Can't use oauth if not a verified email.
		}
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	return
}

var (
	// TODO: Some random string, random for each request
	oauthStateString = "random"
)

type loginCallbackData struct {
	TokenName   string
	Token       string
	RedirectURL string
}

var loginCallbackPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Logging in...</title>
  </head>
  <body>
	<h3>Logging in now...</h3>
  	<script>
		localStorage.setItem('{{.TokenName}}', '{{.Token}}');
		location = '{{.RedirectURL}}'
	</script>
  </body>
</html>
`
