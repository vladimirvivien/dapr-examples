package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	daprClient dapr.Client
	appPort    = os.Getenv("APP_PORT")
	pubsub     = os.Getenv("ORDERS_PUBSUB")
	topic      = os.Getenv("ORDERS_PUBSUB_TOPIC")
)

func main() {
	if appPort == "" {
		appPort = "6060"
	}
	if pubsub == "" {
		pubsub = "orders-pubsub"
	}
	if topic == "" {
		topic = "orders"
	}

	// define subscription
	rcvdSub := &common.Subscription{
		PubsubName: pubsub,
		Topic:      topic,
		Route:      topic,
	}

	// Create service
	s := daprd.NewService(fmt.Sprintf(":%s", appPort))

	// Register handler to handle received orders
	log.Printf("orderprog: registering event handler: {%#v}", rcvdSub)
	if err := s.AddTopicEventHandler(rcvdSub, subHandler); err != nil {
		log.Fatalf("orderproc: topic subscription: %v", err)
	}

	// Set up Dapr client (seems to be a bug that requires client to be created after a service is declared)
	dc, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("order proc: dapr client: %s", err)
	}
	daprClient = dc
	defer daprClient.Close()

	// Start service component last
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("orderproc: starting: %v", err)
	}
}

func subHandler(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
	log.Printf("Subscriber received: %#v", event)
	return false, nil
}
