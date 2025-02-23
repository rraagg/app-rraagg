package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rraagg/rraagg/pkg/routes"
	"github.com/rraagg/rraagg/pkg/services"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Start a new container
	c := services.NewContainer()
	defer func() {
		if err := c.Shutdown(); err != nil {
			c.Web.Logger.Fatal(err)
		}
	}()

	// Build the router
	routes.BuildRouter(c)

	// Start the server
	go func() {
		srv := http.Server{
			Addr:         fmt.Sprintf("%s:%d", c.Config.HTTP.Hostname, c.Config.HTTP.Port),
			Handler:      c.Web,
			ReadTimeout:  c.Config.HTTP.ReadTimeout,
			WriteTimeout: c.Config.HTTP.WriteTimeout,
			IdleTimeout:  c.Config.HTTP.IdleTimeout,
		}

		if c.Config.HTTP.TLS.Enabled {
			autoTLSManager := autocert.Manager{
				Prompt: autocert.AcceptTOS,
				// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
				Cache:      autocert.DirCache("/var/www/.cache"),
				HostPolicy: autocert.HostWhitelist(c.Config.HTTP.TLS.WhiteList),
			}

			srv.TLSConfig = &tls.Config{
				GetCertificate: autoTLSManager.GetCertificate,
			}
		}

		if err := c.Web.StartServer(&srv); err != http.ErrServerClosed {
			c.Web.Logger.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Start the scheduler service to queue periodic tasks
	go func() {
		if err := c.Tasks.StartScheduler(); err != nil {
			c.Web.Logger.Fatalf("scheduler shutdown: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Web.Shutdown(ctx); err != nil {
		c.Web.Logger.Fatal(err)
	}
}
