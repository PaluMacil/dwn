package webserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/api/groupapi"
	"github.com/PaluMacil/dwn/api/infoapi"
	"github.com/PaluMacil/dwn/api/searchapi"
	"github.com/PaluMacil/dwn/api/typeaheadapi"
	"github.com/PaluMacil/dwn/api/userapi"
	"github.com/PaluMacil/dwn/api/usergroupapi"
	"github.com/PaluMacil/dwn/auth"
	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/sections/shopping/shoppingapi"
	"github.com/PaluMacil/dwn/spa"
)

type WebServer struct {
	configuration.WebServerConfiguration
	mux *http.ServeMux
}

func New(db *database.Database, config configuration.Configuration) *WebServer {
	spaModule := spa.New(config.WebServer.ContentRoot)
	authModule := auth.New(db, config)

	userModule := userapi.New(db)
	groupModule := groupapi.New(db)
	usergroupModule := usergroupapi.New(db)
	infoModule := infoapi.New(db, config)
	searchModule := searchapi.New(db)
	typeaheadModule := typeaheadapi.New(db)

	shoppingModule := shoppingapi.New(db)

	ws := &WebServer{
		WebServerConfiguration: config.WebServer,
		mux: http.NewServeMux(),
	}
	ws.mux.Handle("/", spaModule)
	ws.mux.Handle("/oauth/", authModule)
	ws.mux.Handle("/api/user/", userModule)
	ws.mux.Handle("/api/group/", groupModule)
	ws.mux.Handle("/api/usergroup/", usergroupModule)
	ws.mux.Handle("/api/info/", infoModule)
	ws.mux.Handle("/api/search/", searchModule)
	ws.mux.Handle("/api/typeahead/", typeaheadModule)

	ws.mux.Handle("/api/shopping/", shoppingModule)

	return ws
}

func (ws *WebServer) Serve() {
	srv := &http.Server{
		Addr:    ":" + ws.Port,
		Handler: ws.mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Println("Now serving on port", ws.Port)
		log.Println(srv.ListenAndServe())
	}()
	<-stop

	log.Printf("shutting down ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
