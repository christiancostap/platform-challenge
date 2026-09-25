#!/bin/bash
set -euo pipefail

### Expects the image name as argument

if [ "$#" -lt 1 ]; then
    echo "Uso: $0 <image-name>"
    echo "Exemplo: $0 christiancostap/lab"
    exit 1
fi

IMAGE_NAME="$1"
APP_DIR="$(cd "$(dirname "$0")/../app" && pwd)"

### Running Unit Tests

echo "Executing unit tests..."
cd "$APP_DIR"
go test ./... -v

### Generating new hex tag for docker image

NEW_TAG="$(openssl rand -hex 6)"
if [ -z "$NEW_TAG" ]; then
    echo "Error generating tag. Exiting."
    exit 1
fi

IMAGE_TAG="${IMAGE_NAME}:${NEW_TAG}"
echo "New tag: $NEW_TAG"

### Building and pushing docker image

echo "Building image: $IMAGE_TAG"
docker build -t "$IMAGE_TAG" .

echo "Pushing image to registry..."
docker push "$IMAGE_TAG"

echo "Image pushed successfully:"
echo "Image: $IMAGE_NAME"
echo "Tag: $NEW_TAG"

