package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

func main() {
	// Initialize service with Dapr SDK
	svc := daprd.NewService(":8080")

	// Register a route and a handler for service
	if err := svc.AddServiceInvocationHandler("/search", orderHandler); err != nil {
		log.Fatalf("Handler registration: path: /search: %s", err)
	}

	log.Println("Http service started: port 8080")

	// Start service on specified port
	if err := svc.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Service start: %s", err)
	}
}

func orderHandler(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	if in == nil {
		return nil, fmt.Errorf("service handler: parameter missing")
	}

	if len(in.Data) == 0 {
		return nil, fmt.Errorf("/search: invalid search parameter")
	}

	log.Printf("/search: data-type-url: %s, data: %s", in.DataTypeURL, string(in.Data))

	return &common.Content{
		Data:        in.Data,
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}, nil
}
