package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pivotal-cf-experimental/lattice-app/handlers"
	"github.com/pivotal-cf-experimental/lattice-app/helpers"
	"github.com/pivotal-cf-experimental/lattice-app/routes"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/rata"
)

func main() {
	logger := lager.NewLogger("lattice-app")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("lattice-app.starting", lager.Data{"port": port})
	handler, err := rata.NewRouter(routes.Routes, handlers.New(logger))
	if err != nil {
		logger.Fatal("router.creation.failed", err)
	}

	index, err := helpers.FetchIndex()
	go func() {
		t := time.NewTicker(time.Second)
		for {
			<-t.C
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to fetch index: %s\n", err.Error())
			} else {
				fmt.Printf("This is Lattice-App on index: %d\n", index)
			}
		}
	}()

	server := ifrit.Envoke(http_server.New(":"+port, handler))
	logger.Info("lattice-app.up", lager.Data{"port": port})
	err = <-server.Wait()
	if err != nil {
		logger.Error("farewell", err)
	}
	logger.Info("farewell")
}
