#!/bin/bash

# Build the Docker image for arm64 platform
docker buildx build --platform linux/arm64 -t waziupiot/majiup:latest . --no-cache --load

# Push the Docker image to the repository
docker push waziupiot/majiup:latest
