package main

import (
	_ "net/http/pprof" // Import for pprof to register its handlers
	"os"
	"os/signal"
	"syscall"
	"user-service/app"
)

func main() {
	pprofEnabled := os.Getenv("ENABLE_PPROF") == "true"
	service := app.New(app.WithPprofEnabled(pprofEnabled))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	service.Start()

	<-stop
	service.Stop()
}
