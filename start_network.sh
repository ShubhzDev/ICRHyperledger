set -ex

. ./createChannel.sh


# Bring the test network down
docker-compose -f docker-compose.yaml up -d

sleep 10

#
setup_channel
