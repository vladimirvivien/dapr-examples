package main

import (
	"context"
	"log"
	"os"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/google/uuid"
)

var (
	appPort = os.Getenv("APP_PORT")
)

func main() {
	if appPort == "" {
		appPort = "8080"
	}

	dapr := daprd.NewService(":8080")

	// Define service at endpoint /genid
	if err := dapr.AddServiceInvocationHandler("/genid", generateId); err != nil {
		log.Fatalf("genid: invocation handler setup: %v", err)
	}

}

func generateId(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	id := uuid.New()
	out := &common.Content{
		Data:        []byte(id.String()),
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}

	return out, nil
}
