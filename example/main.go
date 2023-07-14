package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Silicon-Ally/rules_oapi_codegen/example/api"
	"github.com/Silicon-Ally/rules_oapi_codegen/example/server"
	"github.com/go-chi/chi/v5"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("no args given")
	}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port = fs.Int("port", 8080, "Port for test HTTP server")
	)
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to load swagger spec: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	petStore := &server.Server{}

	petStoreStrictHandler := api.NewStrictHandler(petStore, nil /* middleware */)

	r := chi.NewRouter()
	r.Use(oapimiddleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	api.HandlerFromMux(petStoreStrictHandler, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%d", *port),
	}

	// And we serve HTTP until the world ends.
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("error running HTTP server: %w", err)
	}

	return nil
}
