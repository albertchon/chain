# echo "y" | bandcli tx multi-send 2000000uband $(bandcli keys list --keyring-backend test | jq -r '.[].address') --from validator --keyring-backend test

./scripts/deploy-yoda.sh 2 test-sentry &
./scripts/deploy-yoda.sh 3 test-validator3 &
./scripts/deploy-yoda.sh 4 test-validator4 &
./scripts/deploy-yoda.sh 5 test-validator5 &
./scripts/deploy-yoda.sh 6 test-validator6 &
./scripts/deploy-yoda.sh 7 test-validator7 &
./scripts/deploy-yoda.sh 8 test-validator8 &
./scripts/deploy-yoda.sh 9 test-validator9 &
./scripts/deploy-yoda.sh 10 test-validator10 &
./scripts/deploy-yoda.sh 11 test-validator11 &
./scripts/deploy-yoda.sh 12 test-validator12 &
./scripts/deploy-yoda.sh 13 test-validator13 &
./scripts/deploy-yoda.sh 14 test-validator14 &
./scripts/deploy-yoda.sh 15 test-validator15 &
./scripts/deploy-yoda.sh 16 test-validator16 &
./scripts/deploy-yoda.sh 17 test-firer &

wait