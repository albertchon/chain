ID=$1
VALIDATOR=validator$ID
NODE=$2 # test-sentry

echo "ID=$ID" > .env
echo "VALIDATOR=$VALIDATOR" >> .env

ssh $NODE <<EOF
sudo systemctl stop bandd
sudo systemctl stop yoda
rm -rf chain*
EOF

scp ../chain.tar.gz $NODE:~/
scp .env $NODE:~/

ssh $NODE <<EOF

echo "export PATH=\$PATH:/usr/local/go/bin:~/go/bin" >> /home/panu/.profile
source ~/.profile

tar -xvf chain.tar.gz
cd chain

make install
# cp ~/go/bin/* /usr/local/bin

./scripts/generate_genesis.sh

sed -E -i \
  's/persistent_peers = \".*\"/persistent_peers = \"0851086afcd835d5a6fb0ffbf96fcdf74fec742e@35.247.144.166:26656,447ff3a50e3141a588df179e760f30ad4720f33d@34.122.233.232:26656\"/' \
  /home/panu/.bandd/config/config.toml

cp genesis.json ~/.bandd/config/
cp -r files ~/.bandd

# prepare
./scripts/start_bandd_with_id.sh $ID
# ./scripts/start_yoda.sh $VALIDATOR

sudo cp bandd.service /etc/systemd/system/bandd.service
# sudo cp yoda.service /etc/systemd/system/yoda.service

sudo systemctl daemon-reload
sudo systemctl start bandd

EOF
