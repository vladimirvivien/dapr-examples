package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/vladimirvivien/daprexamples/types"
)

var (
	daprClient dapr.Client
	appPort    = os.Getenv("APP_PORT")
	stateStore = os.Getenv("ORDERS_STORE")
	pubsub     = os.Getenv("ORDERS_PUBSUB")
	topic      = os.Getenv("ORDERS_PUBSUB_TOPIC")
	route      = os.Getenv("ORDERS_PUBSUB_ROUTE")
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
	if stateStore == "" {
		stateStore = "orders-store"
	}

	// define subscription
	rcvdSub := &common.Subscription{
		PubsubName: pubsub,
		Topic:      topic,
		Route:      route,
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
	orderID, ok := event.Data.(string)
	if !ok {
		err = fmt.Errorf("orders-pubsub: event: unexpected data type: %T", event.Data)
		log.Print(err.Error())
		return false, err
	}

	log.Printf("orders-pubsub: event received: orderID: %s", orderID)

	// retrieve and update order
	orderItem, err := daprClient.GetState(ctx, stateStore, orderID, nil)
	if err != nil {
		log.Printf("orders-pubsub: getstate: %s", err)
		return true, err
	}

	var order types.Order
	if err := json.Unmarshal(orderItem.Value, &order); err != nil {
		log.Printf("orders-pubsub: unmarshal: %s: %s", err, orderItem.Value)
		return false, err
	}
	order.Completed = true

	// save updated order
	orderData, err := json.Marshal(order)
	if err != nil {
		log.Printf("orders-pubsub: marshal: %s", err)
		return false, err
	}

	if err := daprClient.SaveState(ctx, stateStore, orderID, orderData, nil); err != nil {
		log.Printf("orders-pubsub: save state: %s", err)
		return true, err
	}

	log.Printf("orders-pubsub: order update: id: %s", orderID)

	return false, nil
}
