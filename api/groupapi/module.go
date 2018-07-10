package groupapi

import (
	"net/http"
	"strings"

	"github.com/PaluMacil/dwn/app"
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
	route := strings.Split(r.URL.Path, "/")
	switch endpoint := route[1]; endpoint {
	case "":
	}
}
