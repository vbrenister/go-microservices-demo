package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vbrenister/apicommon"
)

const webPort = "80"

type Config struct {
	apicommon.ServerConfig
}

func main() {
	app := Config{}

	log.Printf("Starting broker services at port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Panic(srv.ListenAndServe())

}
