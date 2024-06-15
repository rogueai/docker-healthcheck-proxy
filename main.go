package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rogueai/docker-healthcheck-proxy/healthcheck"
)

var listenPort string = ":3333"

func main() {
	http.HandleFunc("/healthcheck", healthcheck.GetHealthCheck)
	fmt.Printf("Starting HTTP Server: %s", listenPort)
	// instantiate the listener
	var err error = http.ListenAndServe(listenPort, nil)

	// error handling
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	}

	// stops if error
	if err != nil {
		log.Fatalf("error starting server: %s\n", err)
	}

}
