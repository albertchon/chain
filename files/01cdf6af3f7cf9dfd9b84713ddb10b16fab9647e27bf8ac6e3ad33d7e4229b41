#!/usr/bin/env python3

import json
import urllib.request
import sys

COINS_URL = "https://api.coingecko.com/api/v3/coins/list"
VOLUME_URL = (
    "https://api.coingecko.com/api/v3/coins/{}/market_chart?vs_currency=usd&days=1"
)


def make_json_request(url):
    return json.loads(urllib.request.urlopen(url).read())


def main(symbol):
    coins = make_json_request(COINS_URL)
    for coin in coins:
        if coin["symbol"] == symbol.lower():
            slug = coin["id"]
            return make_json_request(VOLUME_URL.format(slug))["total_volumes"][-1][1]
    raise ValueError("unknown CoinGecko symbol: {}".format(symbol))


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
