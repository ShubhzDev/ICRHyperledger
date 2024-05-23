#!/bin/bash

# Stop and remove all Docker containers and networks
docker-compose -f docker-compose.yaml down --volumes --remove-orphans

# Remove any existing channel artifacts
rm -rf channel-artifacts/*.block channel-artifacts/*.tx
