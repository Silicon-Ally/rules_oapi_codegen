package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Silicon-Ally/rules_oapi_codegen/example/api"
	"github.com/Silicon-Ally/rules_oapi_codegen/example/server"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

func TestServerAndClient(t *testing.T) {
	sockPath := startServerOnSocket(t)
	httpC := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockPath)
			},
		},
	}

	// The host (i.e. `localhost`) here isn't relevant, since we override
	// DialContext in our *http.Client. We just need the scheme to tell the client
	// it should use HTTP.
	c, err := api.NewClientWithResponses("http://localhost", api.WithHTTPClient(httpC))
	if err != nil {
		t.Fatalf("failed to init Petstore client: %v", err)
	}

	ctx := context.Background()

	checkPets := func(want []api.Pet) {
		got, err := c.FindPetsWithResponse(ctx, &api.FindPetsParams{})
		if err != nil {
			t.Fatalf("failed to find pets: %v", err)
		}
		if got.JSON200 == nil {
			t.Fatalf("non-200 response %d", got.HTTPResponse.StatusCode)
		}

		if diff := cmp.Diff(*got.JSON200, want); diff != "" {
			t.Errorf("unexpected pets returned (-want +got)\n%s", diff)
		}
	}

	checkPets([]api.Pet{})

	addPet := func(req api.AddPetJSONRequestBody) {
		resp, err := c.AddPetWithResponse(ctx, req)
		if err != nil {
			t.Fatalf("failed to add pet: %v", err)
		}
		if resp.JSON200 == nil {
			t.Fatalf("non-200 response %d", resp.HTTPResponse.StatusCode)
		}
	}

	addPet(api.AddPetJSONRequestBody{
		Name: "Spike",
		Tag:  ptr("good boy"),
	})

	checkPets([]api.Pet{
		{Id: 1, Name: "Spike", Tag: ptr("good boy")},
	})

	addPet(api.AddPetJSONRequestBody{
		Name: "Spot",
		Tag:  ptr("good girl"),
	})

	checkPets([]api.Pet{
		{Id: 1, Name: "Spike", Tag: ptr("good boy")},
		{Id: 2, Name: "Spot", Tag: ptr("good girl")},
	})

	resp, err := c.DeletePetWithResponse(ctx, 1)
	if err != nil {
		t.Fatalf("failed to delete pet: %v", err)
	}
	if resp.JSONDefault != nil {
		t.Fatalf("non-200 response %d", resp.HTTPResponse.StatusCode)
	}

	checkPets([]api.Pet{
		{Id: 2, Name: "Spot", Tag: ptr("good girl")},
	})
}

func ptr[T any](in T) *T {
	return &in
}

func startServerOnSocket(t *testing.T) string {
	swagger, err := api.GetSwagger()
	if err != nil {
		t.Fatalf("failed to load swagger spec: %v", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match.
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
	}

	dir, err := os.MkdirTemp("", "oapitest-sock")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to remove temp test dir %q: %v", dir, err)
		}
	})

	sockPath := filepath.Join(dir, "server.sock")
	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to listen on socket: %v", err)
	}

	go s.Serve(lis)

	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Logf("failed to close server: %v", err)
		}
	})

	return sockPath
}
