package spa

import (
	"net/http"
	"path"
	"strings"
)

type Config struct {
	Path    string
	Project string
}

func (c Config) String() string {
	return path.Join(c.Path, c.Project)
}

func (c Config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	root := c.String()
	fs := http.FileServer(http.Dir(root))
	pathParts := strings.Split(r.URL.Path, `/`)
	lastPart := pathParts[len(pathParts)-1]
	if strings.Contains(lastPart, ".") {
		fs.ServeHTTP(w, r)
	} else {
		http.ServeFile(w, r, path.Join(root, "index.html"))
	}
}
