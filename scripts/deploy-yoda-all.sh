# echo "y" | bandcli tx multi-send 2000000uband $(bandcli keys list --keyring-backend test | jq -r '.[].address') --from validator --keyring-backend test

./scripts/deploy-yoda.sh 2 test-sentry &
./scripts/deploy-yoda.sh 3 panu@35.186.146.40 &
./scripts/deploy-yoda.sh 4 panu@34.87.157.46 &
./scripts/deploy-yoda.sh 5 panu@35.247.144.166 &
./scripts/deploy-yoda.sh 6 panu@35.240.186.168 &
./scripts/deploy-yoda.sh 7 panu@35.198.234.27 &
./scripts/deploy-yoda.sh 8 panu@35.197.155.123 &
./scripts/deploy-yoda.sh 9 panu@35.247.178.105 &
./scripts/deploy-yoda.sh 10 panu@34.87.140.17 &
./scripts/deploy-yoda.sh 11 panu@34.126.117.93 &
./scripts/deploy-yoda.sh 12 panu@34.126.90.149 &
./scripts/deploy-yoda.sh 13 panu@35.240.192.40 &
./scripts/deploy-yoda.sh 14 panu@34.87.115.122 &
./scripts/deploy-yoda.sh 15 panu@35.240.241.225 &
./scripts/deploy-yoda.sh 16 panu@34.87.143.228 &
./scripts/deploy-yoda.sh 17 panu@34.87.115.167 &

wait