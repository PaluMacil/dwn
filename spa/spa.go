package spa

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/PaluMacil/dwn/app"
)

type Module struct {
	*app.App
	ContentRoot string
}

func New(app *app.App) *Module {
	return &Module{
		App:         app,
		ContentRoot: os.Getenv("DWN_CONTENT_ROOT"),
	}
}

func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(mod.ContentRoot))
	pathParts := strings.Split(r.URL.Path, `/`)
	lastPart := pathParts[len(pathParts)-1]
	if strings.Contains(lastPart, ".") {
		fs.ServeHTTP(w, r)
	} else {
		http.ServeFile(w, r, path.Join(mod.ContentRoot, "index.html"))
	}
}
