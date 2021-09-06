package spa

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type Config struct {
	Path    string
	Project string
}

var initOnce sync.Once

func New(path, project string) (Config, error) {
	var initError error = nil
	initOnce.Do(func() {
		var allRoots []string
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			contentRoot := "none specified"
			if path != "" {
				contentRoot = path
			}
			initError = fmt.Errorf("reading the contents of the content root %s: %w", contentRoot, err)
			return
		}
		for _, dir := range dirEntries {
			if dir.IsDir() {
				allRoots = append(allRoots, dir.Name())
			}
		}
		if len(allRoots) == 0 {
			initError = fmt.Errorf("no projects in content root")
			return
		}
		log.Printf("projects in content root %s: %v\n", path, allRoots)
	})
	if initError != nil {
		return Config{}, fmt.Errorf("during initOnce(): %w", initError)
	}

	config := Config{
		Path:    path,
		Project: project,
	}
	if _, err := os.Stat(config.String()); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("folder for project %s does not exist: %w", config.Project, err)
	}
	log.Printf("serving SPA dwn-ui from %s\n", config)
	if _, err := os.Stat(config.Index()); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("index.html for project %s does not exist: %w", config.Project, err)
	}

	return config, nil
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
		http.ServeFile(w, r, c.Index())
	}
}

func (c Config) Index() string {
	return path.Join(c.String(), "index.html")
}
