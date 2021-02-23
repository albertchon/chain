echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator --recover --keyring-backend test --account 0
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester --recover --keyring-backend test --account 0

seq 1 15 | while read i; do
    echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator$i --recover --keyring-backend test --account $i;
    bandd add-genesis-account validator$i 10000000000000uband --keyring-backend test;
done

seq 1 15 | while read i; do
    echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester$i --recover --keyring-backend test --account $i;
bandd add-genesis-account requester$i 10000000000000uband --keyring-backend test;
done

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true