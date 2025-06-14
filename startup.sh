#!/bin/bash

echo "tidying modules..."

go mod tidy

echo "tidy complete!"

echo "starting development server with hot-reloading..."

air -c .air.toml
