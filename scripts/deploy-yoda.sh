ID=$1
VALIDATOR=validator$ID
NODE=$2 # test-sentry

echo "ID=$ID" > .env
echo "VALIDATOR=$VALIDATOR" >> .env


sleep 5

# ssh $NODE <<EOF
# sudo systemctl stop yoda
# EOF

# scp ../chain.tar.gz $NODE:~/
# scp .env $NODE:~/

ssh $NODE <<EOF

cd chain
sudo systemctl stop yoda

echo "export PATH=\$PATH:/usr/local/go/bin:~/go/bin" >> /home/panu/.profile
source ~/.profile

# prepare
# ./scripts/start_bandd_with_id.sh $ID
./scripts/start_yoda.sh $VALIDATOR 8

# sudo cp bandd.service /etc/systemd/system/bandd.service
sudo cp yoda.service /etc/systemd/system/yoda.service

sudo systemctl daemon-reload
sudo systemctl start yoda

EOF
