package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/ec-tmpl/internal/api"
	"example.com/ec-tmpl/internal/config"
	"example.com/ec-tmpl/internal/dependencies"
)

func main() {
	settings, err := config.Load()
	if err != nil {
		log.Printf("settings error: %v", err)
		os.Exit(1)
	}

	greetingService, err := dependencies.BuildGreetingService(settings)
	if err != nil {
		log.Printf("dependencies error: %v", err)
		os.Exit(1)
	}

	server := api.NewHTTPServer(api.Dependencies{
		GreetingService: greetingService,
	})

	addr := fmt.Sprintf("%s:%d", settings.BindIP, settings.BindPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(settings.ShutdownGraceSec)*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("listening on %s", addr)
	if err := server.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
