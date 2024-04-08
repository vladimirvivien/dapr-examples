package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	appPort    = os.Getenv("APP_PORT")
	stateStore = "orders-store"
	genidsvcId = "genidsvc"

	daprClient dapr.Client
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

	dc, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("dapr client: NewClient: %s", err)
	}
	daprClient = dc
	defer daprClient.Close()

	log.Printf("frontend: starting service: port %s", appPort)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders/new", postOrder)
	mux.HandleFunc("GET /orders/order/{id}", getOrder)

	if err := http.ListenAndServe(":"+appPort, mux); err != nil {
		log.Fatalf("frontend: %s", err)
	}
}

func postOrder(w http.ResponseWriter, r *http.Request) {
	var receivedOrder Order
	if err := json.NewDecoder(r.Body).Decode(&receivedOrder); err != nil {
		log.Printf("order decoder: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}

	// invoke genidsvc service to generate order UUID
	out, err := daprClient.InvokeMethod(r.Context(), genidsvcId, "genid", "post")
	if err != nil {
		log.Printf("order genid: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}
	orderID := fmt.Sprintf("order-%s", string(out))
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

	if err := daprClient.SaveState(r.Context(), stateStore, orderID, orderData, nil); err != nil {
		log.Printf("dapr save state: %s", err)
		http.Error(w, "unable to post order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"order":"%s", "status":"received"}`, orderID)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	data, err := daprClient.GetState(r.Context(), stateStore, id, nil)
	if err != nil {
		log.Printf("get order data: %s", err)
		http.Error(w, "unable to get order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(data.Value))

}
