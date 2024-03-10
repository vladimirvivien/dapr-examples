package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/http"
)

const (
	pubsubSvc    = "pubsub-orders"
	pubsubTopic  = "topic-orders"
	stateStore   = "store-orders"
	stateBinding = "binding-orders"
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
		port = "8080"
	}

	service := http.NewService(fmt.Sprintf(":%s", port))

	// Register handler for /orders/new endpoint
	service.AddServiceInvocationHandler("/orders", handleOrder)

	// Start service
	log.Print("Starting order service on port 8080...")
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

	// Augment order data
	orderID := "order-" + fmt.Sprintf("%x", rand.Int31())
	log.Println("Creating order: ", orderID)
	receivedOrder.ID = orderID
	receivedOrder.Completed = false

	// marshal order for downstream processing
	orderData, err := json.Marshal(receivedOrder)
	if err != nil {
		return nil, fmt.Errorf("/orders/new: encode new order: %s", err)
	}
	// Save order to state store
	daprClient, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("error starting service: %s", err)
	}
	defer daprClient.Close()

	log.Printf("Created Dapr client: %#v ", daprClient)

	if err := daprClient.InvokeOutputBinding(ctx, &dapr.InvokeBindingRequest{
		Name:      stateBinding,
		Operation: "create",
		Data:      orderData,
		Metadata:  map[string]string{"key": orderID},
	}); err != nil {
		return nil, fmt.Errorf("/orders/new: save order: %s: %s", orderID, err)
	}

	// Publish order ID to topic
	// if err := daprClient.PublishEvent(ctx, pubsubSvc, pubsubTopic, orderID); err != nil {
	// 	return nil, fmt.Errorf("/orders/new: publish order: %s", err)
	// }

	// return result
	return &common.Content{
		ContentType: "application/json",
		Data:        []byte(fmt.Sprintf(`{"order":"%s", "status":"received"}`, orderID)),
	}, nil
}
