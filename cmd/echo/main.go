package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	echoPort := os.Getenv("ECHO_PORT")
	if echoPort == "" {
		// set default echoPort
		echoPort = "9094"
	}
	historyPort := os.Getenv("ECHO_HISTORY_PORT")
	if historyPort == "" {
		// set default historyPort
		historyPort = "9095"
	}

	history := make(History, 0, 102)

	mux := http.NewServeMux()
	mux.Handle("/", EchoModule{&history})
	srv := &http.Server{
		Addr:    ":" + echoPort,
		Handler: mux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	historyMux := http.NewServeMux()
	historyMux.Handle("/", HistoryModule{&history})
	srvHistory := &http.Server{
		Addr:    ":" + historyPort,
		Handler: historyMux,
		//time from when the connection is accepted to when the request body is fully read
		ReadTimeout: 5 * time.Second,
		//time from the end of the request header read to the end of the response write
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Println("now serving on echoPort", echoPort)
		log.Println(srv.ListenAndServe())
	}()
	go func() {
		log.Println("now serving on historyPort", historyPort)
		log.Println(srvHistory.ListenAndServe())
	}()
	<-stop

	log.Printf("shutting down ...\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	if err := srvHistory.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
