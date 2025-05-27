package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"rest-app/cmd/rest"
	"rest-app/config"
	appSetup "rest-app/internal/setup"
)

func main() {
	// Initialize configuration
	config.InitConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// App setup
	setup := appSetup.Init()

	// Start REST service
	httpServerInstance := rest.StartServer(setup)

	// gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-quit

	log.Println("Shutting down services...")

	// Stop REST server
	if err := httpServerInstance.Shutdown(ctx); err != nil {
		log.Println("Error shutting down server:", err)
	}

	log.Println("Server exited gracefully")
}
