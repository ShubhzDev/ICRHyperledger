export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

sudo chmod +x startNetwork.sh stopNetwork.sh generateArtifacts.sh

./generateArtifacts.sh
./startNetwork.sh
./stopNetwork.sh

export FABRIC_CFG_PATH=$PWD
