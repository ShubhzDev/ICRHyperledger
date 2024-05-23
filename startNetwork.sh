#!/bin/bash

# Start the Docker containers for the network
docker-compose -f docker-compose.yaml up -d

# Wait for a few seconds for the containers to start
sleep 10

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
