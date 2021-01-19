package echo

import (
	"encoding/json"
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"net/http"
	"sort"
	"sync"
	"time"
)

type History []Record

var historyLock = &sync.Mutex{}

func (h *History) add(request string) {
	history := *h
	record := Record{request, time.Now()}
	historyLock.Lock()
	if len(history) >= 100 {
		_, history = history[0], history[1:]
	}
	*h = append(history, record)
	historyLock.Unlock()
}
func (h History) write(w http.ResponseWriter) error {
	historyLock.Lock()
	defer historyLock.Unlock()
	sort.Slice(h, func(i, j int) bool {
		return h[i].Date.UnixNano() > h[j].Date.UnixNano()
	})
	err := json.NewEncoder(w).Encode(h)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
	}
	return nil
}

type Record struct {
	Request string    `json:"request"`
	Date    time.Time `json:"date"`
}

type HistoryModule struct {
	history *History
}

func (mod HistoryModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return
	} else if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	err := mod.history.write(w)
	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
	}
}

func (mod HistoryModule) historyHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if err := cur.Can(core.PermissionViewEchoHistory); err != nil {
		return err
	}
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		return nil
	} else if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return nil
	}
	return mod.history.write(w)
}
