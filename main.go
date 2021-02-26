package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/oguzhn/KeyValueStore/application"
	"github.com/oguzhn/KeyValueStore/cache"
	"github.com/oguzhn/KeyValueStore/controller"
	"github.com/oguzhn/KeyValueStore/filewriter"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	cache := cache.NewCache()

	filewriter := filewriter.NewFileWriter("/tmp/%s-db.txt")

	app := application.NewApplication(cache, filewriter, ctx)

	restController := controller.NewController(app)

	server := &http.Server{
		Addr:    "localhost:9001",
		Handler: restController,
	}

	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()
		log.Println("Server started to listen")
		err := server.ListenAndServe()
		if err != nil {
			log.Println("server erred: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		err := app.WriteToFile()
		if err != nil {
			log.Println("could not write to file: ", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case <-sigCh:
		cancel()
		log.Println("os signal caught ending")
	case <-ctx.Done():
		log.Println("context done")
	}

	sctx, sctxcancel := context.WithTimeout(context.Background(), time.Second)
	defer sctxcancel()
	err := server.Shutdown(sctx)
	if err != nil {
		log.Println("error closing server:", err)
	}
	wg.Wait()
}
