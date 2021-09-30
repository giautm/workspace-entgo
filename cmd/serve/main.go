package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"viecco.dev/awesome/ent"
	elk "viecco.dev/awesome/ent/http"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// Create the ent client.
	c, err := ent.Open("sqlite3", "./ent.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()

	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	// Router and Logger.
	r := chi.NewRouter().With(sentryHandler.Handle)
	l := zap.NewExample()

	// Create the user handler.
	r.Route("/v1", elk.NewUserHandler(c, l).MountRoutes)

	// Start listen to incoming requests.
	fmt.Println("Server running")
	defer fmt.Println("Server stopped")

	sentry.CaptureMessage("It works!")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
