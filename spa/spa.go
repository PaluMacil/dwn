package spa

import (
	"net/http"
	"path"
	"strings"
)

type Module struct {
	contentRoot string
}

func New(contentRoot string) *Module {
	return &Module{
		contentRoot: contentRoot,
	}
}

func (mod Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir(mod.contentRoot))
	pathParts := strings.Split(r.URL.Path, `/`)
	lastPart := pathParts[len(pathParts)-1]
	if strings.Contains(lastPart, ".") {
		fs.ServeHTTP(w, r)
	} else {
		http.ServeFile(w, r, path.Join(mod.contentRoot, "index.html"))
	}
}
