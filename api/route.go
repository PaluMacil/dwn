package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/database"
)

func routeSplitter(c rune) bool {
	return c == '/'
}

func GetRoute(w http.ResponseWriter, r *http.Request, db *database.Database) Route {
	segments := strings.FieldsFunc(r.URL.Path, routeSplitter)
	cur, _ := auth.GetCurrent(r, db)
	var s1, s2, s3 string
	if len(segments) >= 2 {
		s1 = segments[1]
	}
	if len(segments) >= 3 {
		s2 = segments[2]
	}
	if len(segments) >= 4 {
		s3 = segments[3]
	}
	return Route{
		W:        w,
		R:        r,
		Name:     s1,
		Endpoint: s2,
		ID:       s3,
		Current:  cur,
		Db:       db,
	}
}

type Route struct {
	W        http.ResponseWriter
	R        *http.Request
	Name     string
	Endpoint string
	ID       string
	Current  *auth.Current
	Db       *database.Database
}

type API interface {
	API() Route
}

func (rt Route) String() string {
	if rt.Current == nil {
		return "route without identity" //TODO: finish this stringer method
	}
	return fmt.Sprintf("%s %s %s", rt.Name, rt.Endpoint, rt.ID)
}
