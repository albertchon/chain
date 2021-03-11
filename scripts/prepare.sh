make install
./scripts/generate_genesis.sh

cp ~/.bandd/config/genesis.json .

rm -rf files
cp -r ~/.bandd/files .

cd ..
tar -zcvf chain.tar.gz chain
cd chain

 sed -i '' \
  's/persistent_peers = \".*\"/persistent_peers = \"0851086afcd835d5a6fb0ffbf96fcdf74fec742e@35.247.144.166:26656,447ff3a50e3141a588df179e760f30ad4720f33d@34.122.233.232:26656\"/' \
  ~/.bandd/config/config.toml

# cp yoda-init.txt /tmp/yoda-init.txt
# /tmp/yoda-init.txt