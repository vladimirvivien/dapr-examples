#! /bin/bash

curl -i -d '{ "item": "automobile"}' \
    -H "Content-type: application/json" \
    "http://localhost:3500/v1.0/invoke/ordersvc/method/orders/new"