package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/http"
)


type Order struct {
	ID        string
	Items     []string
	Completed bool
}

func main() {
	// Create a Dapr service
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	service := http.NewService(fmt.Sprintf(":%s", port))

	// Register handler for /orders/new endpoint
	service.AddServiceInvocationHandler("/neworder", handleOrder)

	// Start service
	log.Printf("Starting order service on port %s...", port)
	if err := service.Start(); err != nil {
		log.Fatalf("error starting service: %s", err)
	}
}

func handleOrder(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	// Decode received order
	var receivedOrder Order
	if err := json.Unmarshal(in.Data, &receivedOrder); err != nil {
		return nil, fmt.Errorf("/orders/new: decode new order: %s", err)
	}

	// Order received
	orderID := "order-" + fmt.Sprintf("%x", rand.Int31())
	order := fmt.Sprintf(`{"order":"%s", "status":"received"}`, orderID)
	log.Printf("/orders: order received: %s", order)

	// return handler result
	return &common.Content{
		ContentType: "application/json",
		Data:        []byte(order),
	}, nil
}
