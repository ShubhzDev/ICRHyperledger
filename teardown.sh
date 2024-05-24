docker kill $(docker ps -q)
docker rm -f $(docker ps -aq)
docker rmi $(docker images "dev-*" -q)

sudo rm -rf channel-artifacts/*.block channel-artifacts/*.tx
sudo rm -rf crypto-config/*
sudo rm -rf channel