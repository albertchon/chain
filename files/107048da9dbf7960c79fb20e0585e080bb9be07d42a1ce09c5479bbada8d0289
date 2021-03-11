#!/usr/bin/env python3
import requests
import sys

URL = "https://asia-southeast2-price-caching.cloudfunctions.net/query-price"
HEADERS = {"Content-Type": "application/json"}


def main(symbols):
    payload = {"source": "cmc", "symbols": symbols}
    r = requests.post(URL, headers=HEADERS, json=payload)
    r.raise_for_status()
    pxs = r.json()

    if len(pxs) != len(symbols):
        raise Exception("PXS_AND_SYMBOL_LEN_NOT_MATCH")

    return ",".join(pxs)


if __name__ == "__main__":
    try:
        print(main(sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
