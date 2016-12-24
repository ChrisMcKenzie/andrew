#!/bin/bash

mkdir -p bin
GOOS=linux GOARCH=amd64 go build -o bin/fulfillment main.go
docker build -t quay.io/chrismckenzie/fulfillment .
