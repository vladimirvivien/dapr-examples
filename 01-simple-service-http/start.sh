#! /bin/bash

dapr run --app-id emojisvc \
         --app-protocol http \
         --app-port 8080 \
         --dapr-http-port 3500 \
         --log-level debug \
         go run ./service/main.go