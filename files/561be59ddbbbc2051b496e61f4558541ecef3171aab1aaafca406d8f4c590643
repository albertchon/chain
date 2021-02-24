#!/usr/bin/env python3

import json
import urllib.request
import sys

ALPHAVANTAGE_URL = (
    "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol={}&apikey={}"
)


def make_json_request(url):
    req = urllib.request.Request(url)
    req.add_header(
        "User-Agent",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
    )
    return json.loads(urllib.request.urlopen(req).read())


def main(symbol, api_key):
    random = make_json_request(ALPHAVANTAGE_URL.format(symbol, api_key))
    return random["Global Quote"]["05. price"]


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
