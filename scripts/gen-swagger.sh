#!/bin/bash
set -e

# Directory where the API code lives
API_DIR="./cmd/jobdoe"

# Output directory for generated docs
OUTPUT_DIR="./api/swagger"

# Check if swag is installed
if ! command -v swag &> /dev/null
then
    echo "swag could not be found. Install it with:"
    echo "  go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

echo "Generating Swagger docs for JobDoe..."

# Create output dir if it doesn't exist
mkdir -p $OUTPUT_DIR

# Run swag init
swag init -g "$API_DIR/main.go" -o "$OUTPUT_DIR" --parseDependency --parseInternal

echo "Swagger docs generated in $OUTPUT_DIR"
