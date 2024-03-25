#!/bin/bash

SWAGGER_FILE="./docs/swagger.json"

# Use sed to replace the host field in swagger.json
sed -i "s/\"host\": \".*\"/\"host\": \"$1\"/" $SWAGGER_FILE || true

echo "swagger.json updated with host: $1"
