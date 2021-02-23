NAME=$1
DEST=$2

rm -rf ~/.bandd
mkdir -p $DEST

bandd init $NAME --chain-id bandchain
cp ~/.bandd/config/node_key.json $DEST
cp ~/.bandd/config/priv_validator_key.json $DEST

bandd tendermint show-node-id
bandd tendermint show-validator
