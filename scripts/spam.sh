while true; do
    if [ -f "/tmp/stop.spam" ]; then
        exit 0
    fi
    sleep 2;
    yes | bandcli tx oracle request 14 16 10 -c 0000000a0000000347565400000004494f5354000000034b4559000000044c4f4f4d000000034d4554000000034d4647000000034d4c4e000000034d544c000000034d5942000000054e4558584f000000003b9aca00 -m requester  --from $1 --keyring-backend test --chain-id bandchain --gas 1590000;
done