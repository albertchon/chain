./scripts/generate_genesis.sh

cp ~/.bandd/config/genesis.json .

rm -rf files
cp -r ~/.bandd/files .

cd ..
tar -zcvf chain.tar.gz chain
cd chain

 sed -i '' \
  's/persistent_peers = \".*\"/persistent_peers = \"0851086afcd835d5a6fb0ffbf96fcdf74fec742e@104.154.93.0:26656\"/' \
  ~/.bandd/config/config.toml