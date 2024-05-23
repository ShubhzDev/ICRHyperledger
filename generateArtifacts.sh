#!/bin/bash

# Generate cryptographic material
cryptogen generate --config=./crypto-config.yaml

# Generate channel artifacts
configtxgen -profile OrdererGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile mychannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel
configtxgen -profile mychannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
configtxgen -profile mychannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
configtxgen -profile mychannel -outputAnchorPeersUpdate ./channel-artifacts/ParentOrgMSPanchors.tx -channelID mychannel -asOrg ParentOrgMSP
