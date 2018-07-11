package userapi

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/db"
)

type Module struct {
	*app.App
}

func New(app *app.App) *Module {
	return &Module{
		App: app,
	}
}

func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cur, _ := auth.GetCurrent(r, mod.Db)

	route := strings.Split(r.URL.Path, "/")
	switch endpoint := route[3]; endpoint {
	case "me":
		if cur == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		groups, err := mod.Db.Groups.GroupsFor(cur.User.Email)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		me := struct {
			auth.Current
			Groups []db.Group
		}{
			*cur,
			groups,
		}
		if err := json.NewEncoder(w).Encode(me); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
}
