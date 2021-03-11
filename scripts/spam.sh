while true; do
    if [ -f "/tmp/stop.spam" ]; then
        exit 0
    fi
    sleep 6;
    ./jobs.sh $1
    # sleep 6 && yes | bandcli tx oracle request 1 16 10 -c 0000000a00000002485400000003564554000000054d494f544100000003534e5800000004434f4d50000000034f4d47000000034d4b5200000004444f4745000000034b534d00000003444742000000003b9aca00 -m requester --from requester$1 --keyring-backend test --chain-id bandchain --gas 1590000 --memo 'x';
done