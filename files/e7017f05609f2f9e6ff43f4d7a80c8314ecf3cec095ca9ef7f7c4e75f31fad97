#!/usr/bin/env python3

import requests
import sys


def main(rpc, to):
    headers = {
        "Content-Type": "application/json",
    }
    data = (
        """{ "jsonrpc": "2.0", "method": "eth_call", "params": [ { "to": "%s", "data": "0x18160ddd" }, "latest" ], "id": 1 }"""
        % (to)
    )
    response = requests.post(
        rpc,
        headers=headers,
        data=data,
    )
    return int(response.json()["result"], 16)


if __name__ == "__main__":
    try:
        print(main(sys.argv[1], sys.argv[2]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
