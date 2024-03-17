package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	appPort    = os.Getenv("APP_PORT")
	stateStore = "orders-store"
)

type Order struct {
	ID        string
	Items     []string
	Completed bool
}

func main() {
	if appPort == "" {
		appPort = "8080"
	}
	log.Printf("frontend: starting service: port %s", appPort)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders/new", postOrder)

	if err := http.ListenAndServe(":"+appPort, mux); err != nil {
		log.Fatalf("frontend: %s", err)
	}
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	daprClient, err := dapr.NewClient()
	if err != nil {
		log.Printf("dapr client: NewClient: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}
	defer daprClient.Close()

	var receivedOrder Order
	if err := json.NewDecoder(r.Body).Decode(&receivedOrder); err != nil {
		log.Printf("order decoder: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}

	orderID := fmt.Sprintf("order-%x", rand.Int31())
	receivedOrder.ID = orderID
	receivedOrder.Completed = true
	log.Printf("order received: [orderid=%s]", orderID)

	// marshal order for downstream processing
	orderData, err := json.Marshal(receivedOrder)
	if err != nil {
		log.Printf("order data: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}

	if err := daprClient.SaveState(context.Background(), stateStore, orderID, orderData, nil); err != nil {
		log.Printf("dapr save state: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"order":"%s", "status":"received"}`, orderID)
}
