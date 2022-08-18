#!/bin/sh
image=$registry/teritori/teritorid:$(git rev-parse --short HEAD)
docker build .. -t $image
docker push $image
