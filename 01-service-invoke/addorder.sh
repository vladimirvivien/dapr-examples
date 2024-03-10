#! /bin/bash

curl -i -d '{ "item": "automobile"}' \
    -H "Content-type: application/json" \
    "http://localhost:8080/orders"