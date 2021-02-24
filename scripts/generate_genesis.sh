DIR=`dirname "$0"`

rm -rf ~/.band*

# initial new node
bandd init validator --chain-id bandchain
echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator --recover --keyring-backend test --account 0
echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester --recover --keyring-backend test --account 0

seq 1 17 | while read i; do
    echo "lock nasty suffer dirt dream fine fall deal curtain plate husband sound tower mom crew crawl guard rack snake before fragile course bacon range" \
    | bandcli keys add validator$i --recover --keyring-backend test --account $i;
    bandd add-genesis-account validator$i 10000000000000uband --keyring-backend test;
done

seq 1 50 | while read i; do
    echo "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comic" \
    | bandcli keys add requester$i --recover --keyring-backend test --account $i;
    bandd add-genesis-account requester$i 10000000000000uband --keyring-backend test;
done

# add accounts to genesis
bandd add-genesis-account validator 10000000000000uband --keyring-backend test
bandd add-genesis-account requester 10000000000000uband --keyring-backend test

# genesis configurations
bandcli config chain-id bandchain
bandcli config output json
bandcli config indent true
bandcli config trust-node true

# register initial validators
bandd gentx \
    --amount 100000000uband \
    --node-id 11392b605378063b1c505c0ab123f04bd710d7d7 \
    --pubkey bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg \
    --name validator \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 0851086afcd835d5a6fb0ffbf96fcdf74fec742e \
    --pubkey bandvalconspub1addwnpepqfey4c5ul6m5juz36z0dlk8gyg6jcnyrvxm4werkgkmcerx8fn5g2gj9q6w \
    --name validator2 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 7b58b086dd915a79836eb8bfa956aeb9488d13b0 \
    --pubkey bandvalconspub1addwnpepqwj5l74gfj8j77v8st0gh932s3uyu2yys7n50qf6pptjgwnqu2arxkkn82m \
    --name validator3 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 63808bd64f2ec19acb2a494c8ce8467c595f6fba \
    --pubkey bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj \
    --name validator4 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 3e13c71d4d9eb5e24098e0ec0fde5ba5ba81e1f7 \
    --pubkey bandvalconspub1addwnpepqw6mfpx4lyqt2wv3fvrmur0p93fv094hemkky99gwavh4ve53lc9ct9a5h2 \
    --name validator5 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id ebff6e173cc636ab5d4ac5b1ace03bea01ade1be \
    --pubkey bandvalconspub1addwnpepq06c9a8ngqwm8knsq54n6ntal46rkz03kvaq5tflz9ugzjxkawms6tf6uqf \
    --name validator6 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 9011c9cc52f79ab67b79c23f3ecb47d17c8a7ef6 \
    --pubkey bandvalconspub1addwnpepqtqnctwmzzpx6xvzkl4fuvmxs996fw0uunu9wyzpuv8alyvpjp0acte84ey \
    --name validator7 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 011be81810266ac71aedbb3ef4d0b1c39e28e3b6 \
    --pubkey bandvalconspub1addwnpepq0zf9yzdp6wk9jqwlwjh33tulem583nw2rkfw6tkpfd4xv8642dk7xvv6yc \
    --name validator8 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id d5916cb79894aaccf24c583a13a044652b6a8d2d \
    --pubkey bandvalconspub1addwnpepqv7zzecf9e54yx6c2ceepmzjlasqu05dfrzt4xmfn9xkuzuntg68uza5hkq \
    --name validator9 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id efb6db668774c2d6315d0e56ac51bace05c47af5 \
    --pubkey bandvalconspub1addwnpepq00pzf5hvzgflx3ejxg9e4dwfuwvmcd5x3j5wk0fp8gl3vyukuztq7dpgrx \
    --name validator10 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 025be7a59b14e37693a58e37bc955beab5830764 \
    --pubkey bandvalconspub1addwnpepq2ywu8halv3gahyn5uyjwymsxfcywvr25hmeg56270588qcmvfwqw8t54n2 \
    --name validator11 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 0607b43ef32006189ce66430aa5d44081cc60160 \
    --pubkey bandvalconspub1addwnpepqws0da2faadh33v032ckzynqecw877um66elx4wg87rxnyr2nnq5kmqetrw \
    --name validator12 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id e44470a6a73798bf69f51d47424ece63fbf97b9c \
    --pubkey bandvalconspub1addwnpepqd9pahv0qkhwxvh2xky73p6jd89xwgwq6uj4zpqmpse84758f42rw2z5yly \
    --name validator13 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id c514083635dd6390f79adc97b2248ac8c59f6b9b \
    --pubkey bandvalconspub1addwnpepqf2sjrj8d6675ks8sc8ztc3tcecdgr7rzz9xxfjp3wvl8qx4jks2z9ywr2g \
    --name validator14 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 4d87ff9809c462a26da3a9577b32bfa35b385722 \
    --pubkey bandvalconspub1addwnpepqtyrrnthdulme2sr7lfud9p4qn7fpl0hq6305ztclt4ctr2p5xa32xdalr9 \
    --name validator15 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 17f8baf6a124674180d447b3384da504978079fa \
    --pubkey bandvalconspub1addwnpepq0fw9cykzxyfhlt5n8e0tyrpq3lxt3puguaznpzm8e35d5rgcp4w5kn3hzm \
    --name validator16 \
    --keyring-backend test

bandd gentx \
    --amount 100000000uband \
    --node-id 447ff3a50e3141a588df179e760f30ad4720f33d \
    --pubkey bandvalconspub1addwnpepqgct6nxy6vhxfma2zwffk7cufzpsvy98da8nyaun97hs0jgwud95g9raq4e \
    --name validator17 \
    --keyring-backend test

/Users/thebevrishot/Workspaces/genesis_ds_os/genesis/scripts/add_os_ds.sh

# collect genesis transactions
bandd collect-gentxs


