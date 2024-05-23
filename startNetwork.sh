#!/bin/bash

# Start the Docker containers for the network
docker-compose -f docker-compose.yaml up -d

# Wait for a few seconds for the containers to start
sleep 10

# Check if the required containers are running
if [ "$(docker ps -q -f name=orderer.example.com)" ]; then
    echo "Orderer container is running."
else
    echo "Orderer container is not running. Exiting..."
    exit 1
fi

if [ "$(docker ps -q -f name=peer0.org1.example.com)" ]; then
    echo "Peer for Org1 container is running."
else
    echo "Peer for Org1 container is not running. Exiting..."
    exit 1
fi

if [ "$(docker ps -q -f name=peer0.org2.example.com)" ]; then
    echo "Peer for Org2 container is running."
else
    echo "Peer for Org2 container is not running. Exiting..."
    exit 1
fi

if [ "$(docker ps -q -f name=peer0.parentorg.example.com)" ]; then
    echo "Peer for ParentOrg container is running."
else
    echo "Peer for ParentOrg container is not running. Exiting..."
    exit 1
fi

# Create and join the channel
docker exec cli peer channel create -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/mychannel.tx --outputBlock ./channel-artifacts/mychannel.block
docker exec cli peer channel join -b ./channel-artifacts/mychannel.block

# Anchor peer updates for Org1
docker exec cli peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx

# Anchor peer updates for Org2
docker exec cli peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx

# Anchor peer updates for ParentOrg
docker exec cli peer channel update -o orderer.example.com:7050 -c mychannel -f ./channel-artifacts/ParentOrgMSPanchors.tx

# Install and instantiate chaincode
docker exec cli peer chaincode install -n intercompany -v 1.0 -l golang -p chaincode/intercompany/go
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n intercompany -l golang -v 1.0 -c '{"Args":[]}' -P "OR('Org1MSP.peer','Org2MSP.peer','ParentOrgMSP.peer')"
