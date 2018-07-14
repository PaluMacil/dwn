package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/db"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Module struct {
	*app.App
	OAuth OAuth
}

func New(app *app.App) *Module {
	home, err := url.Parse(app.HomePage())
	if err != nil {
		panic("Cannot parse home URL: " + err.Error())
	}
	googleCallbackURL, err := url.Parse("oauth/google/callback")
	if err != nil {
		panic("Cannot parse google callback URL: " + err.Error())
	}
	googleRedirect := home.ResolveReference(googleCallbackURL)
	return &Module{
		App: app,
		OAuth: OAuth{
			Google: Google{
				Key:    os.Getenv("DWN_OAUTH_GOOGLE_KEY"),
				Secret: os.Getenv("DWN_OAUTH_GOOGLE_SECRET"),
				Config: &oauth2.Config{
					RedirectURL:  googleRedirect.String(),
					ClientID:     os.Getenv("DWN_OAUTH_GOOGLE_KEY"),
					ClientSecret: os.Getenv("DWN_OAUTH_GOOGLE_SECRET"),
					Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
						"https://www.googleapis.com/auth/userinfo.email"},
					Endpoint: google.Endpoint,
				},
			},
		},
	}
}

type OAuth struct {
	Google Google
}

type Google struct {
	Key    string
	Secret string
	Config *oauth2.Config
}

func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/") //route[2]) is google
	if len(route) != 4 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch endpoint := route[3]; endpoint {
	case "login":
		url := mod.OAuth.Google.Config.AuthCodeURL(oauthStateString)
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
		token, err := mod.OAuth.Google.Config.Exchange(oauth2.NoContext, code)
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
			//var session db.Session
			session := mod.Db.Sessions.GenerateFor(claims.Email)
			err := mod.Db.Sessions.Set(session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//if user exists in database, save session, update last login
			user, err := mod.Db.Users.Get(claims.Email)
			if db.IsKeyNotFoundErr(err) {
				user = claims.CreateUser()
				if claims.Email == mod.Setup.InitialAdmin {
					//TODO: handle err below and add other users to User group
					mod.Db.UserGroups.Set(db.UserGroup{
						Email:     claims.Email,
						GroupName: db.BuiltInGroupAdmin,
					})
				}
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user.LastLogin = time.Now()
			err = mod.Db.Users.Set(user)
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
					TokenName:   "dwn-token",
					Token:       session.Token,
					RedirectURL: mod.HomePage(), //TODO: Check to see if a different redirect is requested and if it is safe
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

type Current struct {
	User    db.User    `json:"user"`
	Session db.Session `json:"session"`
	db      *db.Db     `json:"-"`
}

func GetCurrent(r *http.Request, db *db.Db) (*Current, error) {
	token := r.Header.Get("dwn-token")
	if token == "" {
		return nil, nil
	}
	session, err := db.Sessions.Get(token)
	if err != nil {
		return nil, err
	}
	user, err := db.Users.Get(session.Email)
	if err != nil {
		return nil, err
	}
	return &Current{
		User:    user,
		Session: session,
		db:      db,
	}, nil
}

func (c *Current) Can(permission string) (bool, error) {
	groups, err := c.db.Groups.GroupsFor(c.User.Email)
	if err != nil {
		return false, err
	}
	for _, g := range groups {
		if g.HasPermission(permission) {
			return true, nil
		}
	}
	return false, nil
}

func (c *Current) Is(groupName string) (bool, error) {
	return c.db.UserGroups.Exists(c.User.Email, groupName)
}

var (
	// Some random string, random for each request
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
	Logging in now...
  	<script>
		localStorage.setItem('{{.TokenName}}', '{{.Token}}');
		location = '{{.RedirectURL}}'
	</script>
  </body>
</html>
`
