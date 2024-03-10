## Simple service invocation

This directory contains a simple example that shows how to create a service invocation handler using the 
Dapr API.  The service (`ordersvc`) is setup to listen for incoming HTTP POST requests with JSON payload.
For this introductory example, the code simply returns a response acknowledging that the order has been created.

### The service

The `main` function creates the entry point for the service. Notice that the code is using Dapr's `http` package instead of the standard library's. 

```go
package main

import (
    ...
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/http"
)

func main() {
	// Setup service port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

    // Create new service instance
	service := http.NewService(fmt.Sprintf(":%s", port))

	// Register handler for service endpoint
	service.AddServiceInvocationHandler("/orders", handleOrder)

	// Start service
	log.Printf("Starting order service on port %s ...", port)
	if err := service.Start(); err != nil {
		log.Fatalf("error starting service: %s", err)
	}
}
```

### The service handler

In this example, the service handler is a simple function with signature `func (context.Context, *common.InvocationEvent) (*common.Content, error)`. The handler receives incoming service invocation arguments via variable `in` containing everything needed to handle the request.

```go
type Order struct {
	ID        string
	Items     []string
	Completed bool
}
...
func handleOrder(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	// Decode received order
	var receivedOrder Order
	if err := json.Unmarshal(in.Data, &receivedOrder); err != nil {
		return nil, fmt.Errorf("/orders: decode new order: %s", err)
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

```