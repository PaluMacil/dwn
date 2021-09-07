package webserver

import (
	"context"
	"fmt"
	"github.com/PaluMacil/dwn/module/echo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/configuration"
	configurationapi "github.com/PaluMacil/dwn/module/configuration/api"
	coreapi "github.com/PaluMacil/dwn/module/core/api"
	dashboardapi "github.com/PaluMacil/dwn/module/dashboard/api"
	oauthapi "github.com/PaluMacil/dwn/module/oauth/api"
	registrationapi "github.com/PaluMacil/dwn/module/registration/api"
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

func New(db *database.Database) (*WebServer, error) {
	config := db.Config.Get()
	globalMux := mux.NewRouter()
	ws := &WebServer{
		WebServerConfiguration: config.WebServer,
		mux:                    globalMux,
	}

	apiFactory := handler.Factory{
		Db:                    db,
		AssumeJSONContentType: true,
	}
	genericFactory := handler.Factory{
		Db:                    db,
		AssumeJSONContentType: false,
	}

	// Set host subrouters
	var dwnHost, echoHost, echoHistoryHost *mux.Router
	if config.Prod {
		dwnHost = ws.mux.Host("danwolf.net").Subrouter()
		echoHost = ws.mux.Host("echo.danwolf.net").Subrouter()
		echoHistoryHost = ws.mux.Host("echo-history.danwolf.net").Subrouter()
	} else {
		localhostMatchPattern := fmt.Sprintf("{host:localhost:(?:%s|%s)}", ws.Port, ws.UIProxyPort)
		dwnHost = ws.mux.Host(localhostMatchPattern).Subrouter()
		echoHost, echoHistoryHost = dwnHost, dwnHost
	}

	// subdomain subrouters
	echo.RegisterRoutes(echoHistoryHost, apiFactory)
	if config.Prod {
		echoHost.PathPrefix("/").Handler(echo.NewEchoHandler())
	} else {
		echoHost.PathPrefix("/s/echo/").Handler(echo.NewEchoHandler())
		echoHost.Path("/s/echo").Handler(echo.NewEchoHandler())
	}

	// Set module subrouters
	// ...core
	coreRouter := dwnHost.PathPrefix("/api/core/").Subrouter()
	coreapi.RegisterRoutes(coreRouter, apiFactory)
	// ...oauth
	oauthBaseRouter := dwnHost.PathPrefix("/oauth/").Subrouter()
	oauthapi.RegisterBaseRoutes(oauthBaseRouter, genericFactory)
	// ...typeahead
	typeaheadRouter := dwnHost.PathPrefix("/api/typeahead/").Subrouter()
	typeaheadapi.RegisterRoutes(typeaheadRouter, apiFactory)
	// ...shopping
	shoppingRouter := dwnHost.PathPrefix("/api/shopping/").Subrouter()
	shoppingapi.RegisterRoutes(shoppingRouter, apiFactory)
	// ...server
	serverRouter := dwnHost.PathPrefix("/api/server/").Subrouter()
	serverapi.RegisterRoutes(serverRouter, apiFactory)
	// ...configuration
	configurationRouter := dwnHost.PathPrefix("/api/configuration/").Subrouter()
	configurationapi.RegisterRoutes(configurationRouter, apiFactory)
	// ...dashboard
	dashboardRouter := dwnHost.PathPrefix("/api/dashboard/").Subrouter()
	dashboardapi.RegisterRoutes(dashboardRouter, apiFactory)
	// ...registration
	registrationRouter := dwnHost.PathPrefix("/api/registration/").Subrouter()
	registrationapi.RegisterRoutes(registrationRouter, apiFactory)

	// ...content roots
	dwnuiSPA, err := spa.New(config.WebServer.ContentRoot, "dwn-ui")
	if err != nil {
		return nil, fmt.Errorf("creating new SPA for dwn-ui: %w", err)
	}
	dwnHost.PathPrefix("/").Handler(dwnuiSPA)

	// ...error status routes
	globalMux.NotFoundHandler = handler.Status404Handler(config.WebServer.Status404HandlerName)

	return ws, nil
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
	signal.Notify(stop, os.Interrupt, os.Kill)
	go func() {
		log.Println(srv.ListenAndServe())
	}()
	log.Println("Now serving on port", ws.Port)
	<-stop

	log.Printf("shutting down ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
