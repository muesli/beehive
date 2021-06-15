#!/bin/sh

if ! command -v goreleaser; then
  echo "goreleaser not found"
  exit 1
fi

# Get the highest tag number
VERSION="$(git describe --abbrev=0 --tags)"
VERSION=${VERSION:-'0.0.0'}

# Get number parts
MAJOR="${VERSION%%.*}"; VERSION="${VERSION#*.}"
MINOR="${VERSION%%.*}"; VERSION="${VERSION#*.}"
PATCH="${VERSION%%.*}"; VERSION="${VERSION#*.}"

# Increase version
PATCH=$((PATCH+1))

TAG="${1}"

if [ "${TAG}" = "" ]; then
  TAG="${MAJOR}.${MINOR}.${PATCH}"
fi

echo "Releasing ${TAG} ..."

# Tag git and push
git tag -a -s -m "Release ${TAG}" "${TAG}"
git push --tags

# Run goreleaser
goreleaser release --rm-dist

# Build and push Docker image
DOCKER_REPO="fribbledom/beehive"
docker buildx build \
    --progress plain \
    --platform=linux/amd64,linux/386,linux/arm64,linux/arm/v7,linux/arm/v6 \
    -t="$DOCKER_REPO:$TAG" \
    -t="$DOCKER_REPO:latest" \
    --push \
    .
