package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/oauth"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// GET /oauth/google/{step}
func flowHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	config := db.Config.Get()
	switch step := vars["step"]; step {
	case "login":
		url := config.FS.Auth.Google.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	case "callback":
		state := r.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return nil
		}

		code := r.FormValue("code")
		token, err := config.FS.Auth.Google.Exchange(r.Context(), code)
		if err != nil {
			fmt.Printf("Code exchange failed with '%s'\n", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return nil
		}

		response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			return err
		}

		defer response.Body.Close()
		claims := &oauth.GoogleClaims{}
		err = json.NewDecoder(response.Body).Decode(claims)
		if err != nil {
			return err
		}

		if claims.VerifiedEmail {
			//if user exists in database, save session, update last login
			user, err := db.Users.FromEmail(claims.Email)
			if db.IsKeyNotFoundErr(err) {
				displayName, err := generateDisplayName(claims.GivenName)
				if err != nil {
					return err
				}
				id, err := db.NextID()
				if err != nil {
					return err
				}
				user = claims.CreateUser(id, displayName)
				if claims.Email == config.Setup.InitialAdmin {
					//TODO: handle err below and add other users to User group
					err = db.UserGroups.Set(core.UserGroup{
						UserID:    id,
						GroupName: core.BuiltInGroupAdmin,
					})
					if err != nil {
						return err
					}
				}
			} else if err != nil {
				return err
			}
			user.LastLogin = time.Now()
			err = db.Users.Set(user)
			if err != nil {
				return err
			}
			ip := core.IP(r)
			session := db.Sessions.GenerateFor(user.ID, ip)
			err = db.Sessions.Set(session)
			if err != nil {
				return err
			}
			tmpl, err := template.New("login").Parse(loginCallbackPage)
			if err != nil {
				return err
			}
			return tmpl.Execute(w,
				loginCallbackData{
					TokenName:   "dwn-token",
					Token:       session.Token,
					RedirectURL: config.WebServer.HomePage(), //TODO: Check to see if a different redirect is requested and if it is safe
				})
		} else {
			//TODO: Send to registration page. Can't use oauth if not a verified email.
		}
		return nil
	}
	return errs.StatusNotFound
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

const loginCallbackPage = `
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

func generateDisplayName(givenName string) (core.DisplayName, error) {
	//TODO: generate alternatives and check for repeats
	return core.DisplayName(givenName), nil
}
