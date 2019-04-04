package webserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PaluMacil/dwn/configuration"
	"github.com/PaluMacil/dwn/database"
)

type WebServer struct {
	configuration.WebServerConfiguration
	mux *http.ServeMux
}

func New(db *database.Database, config configuration.Configuration) *WebServer {
	ws := &WebServer{
		WebServerConfiguration: config.WebServer,
		mux:                    http.NewServeMux(),
	}

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
