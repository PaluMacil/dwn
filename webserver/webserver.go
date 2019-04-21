package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
	coreapi "github.com/PaluMacil/dwn/module/core/api"
	oauthapi "github.com/PaluMacil/dwn/module/oauth/api"
	serverapi "github.com/PaluMacil/dwn/module/server/api"
	shoppingapi "github.com/PaluMacil/dwn/module/shopping/api"
	typeaheadapi "github.com/PaluMacil/dwn/module/typeahead/api"
	"github.com/PaluMacil/dwn/spa"
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

type WebServer struct {
	configuration.WebServerConfiguration
	mux *mux.Router
}

func New(db *database.Database, config configuration.Configuration) *WebServer {
	ws := &WebServer{
		WebServerConfiguration: config.WebServer,
		mux:                    mux.NewRouter(),
	}

	apiFactory := handler.Factory{
		Db:                    db,
		Config:                config,
		AssumeJSONContentType: true,
	}
	genericFactory := handler.Factory{
		Db:                    db,
		Config:                config,
		AssumeJSONContentType: false,
	}
	// TODO: refactor config and remove this check on the os args
	prod := len(os.Args) == 2 && os.Args[1] == "prod"

	// Set host subrouters
	var dwnHost *mux.Router
	if prod {
		dwnHost = ws.mux.Host("danwolf.net").Subrouter()
	} else {
		localhostMatchPattern := fmt.Sprintf("localhost:(?:%s|%s)", ws.Port, ws.UIProxyPort)
		dwnHost = ws.mux.Host(localhostMatchPattern).Subrouter()
	}

	// Set module subrouters
	// ...core
	coreRouter := dwnHost.PathPrefix("/api/core/").Subrouter()
	coreapi.RegisterRoutes(coreRouter, apiFactory)
	// ...oauth
	oauthBaseRouter := dwnHost.PathPrefix("/oauth/").Subrouter()
	oauthapi.RegisterBaseRoutes(oauthBaseRouter, genericFactory)
	// ...typeahead
	typeaheadRouter := dwnHost.PathPrefix("/typeahead/").Subrouter()
	typeaheadapi.RegisterRoutes(typeaheadRouter, apiFactory)
	// ...shopping
	shoppingRouter := dwnHost.PathPrefix("/shopping/").Subrouter()
	shoppingapi.RegisterRoutes(shoppingRouter, apiFactory)
	// ...server
	serverRouter := dwnHost.PathPrefix("/server/").Subrouter()
	serverapi.RegisterRoutes(serverRouter, apiFactory)
	// ...content roots
	dwnHost.PathPrefix("/").Handler(spa.ContentRoot(config.WebServer.ContentRoot))

	ws.mux.Host("echo.danwolf.net").Subrouter()
	ws.mux.Host("echo-history.danwolf.net").Subrouter()
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
