package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/yusukemisa/registMember/src/handler"

	"go.mercari.io/datastore"

	"github.com/yusukemisa/registMember/src/infra"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	server := &http.Server{
		Addr: "localhost:" + port,
	}

	ctx := context.Background()
	var err error
	datastoreClient, err := infra.NewDatastoreClient(ctx,
		datastore.WithProjectID(os.Getenv("ProjectID")),
		datastore.WithCredentialsFile(os.Getenv("GCP_CREDENTIALS")),
	)

	if err != nil {
		log.Fatal("failed to make new datastore client")
	}

	http.Handle("/user", handler.UserHandler{
		Cli: datastoreClient,
	})
	http.Handle("/", handler.PostHandler{
		Cli: datastoreClient,
	})
	server.ListenAndServe()
}
