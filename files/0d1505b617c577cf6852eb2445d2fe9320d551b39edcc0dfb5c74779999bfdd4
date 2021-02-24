#!/usr/bin/env python3

import json
import urllib.request
import sys

CRYPTOCOMPARE_URL = (
    "https://min-api.cryptocompare.com/data/symbol/histoday?fsym={}&tsym=USD&limit=1"
)


def make_json_request(url):
    return json.loads(urllib.request.urlopen(url).read())


def main(symbol):
    res = make_json_request(CRYPTOCOMPARE_URL.format(symbol))
    return res["Data"][0]["total_volume_total"]


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
