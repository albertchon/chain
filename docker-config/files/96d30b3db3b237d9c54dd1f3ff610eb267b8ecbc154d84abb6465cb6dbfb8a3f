#!/usr/bin/env python3

import json
import urllib.request
import sys

BLOCKCHAIN_INFO_URL = "https://blockchain.info/q/getblockcount"


def make_json_request(url):
    req = urllib.request.Request(url)
    req.add_header(
        "User-Agent",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
    )
    return json.loads(urllib.request.urlopen(req).read())


def main():
    count = make_json_request(BLOCKCHAIN_INFO_URL)
    return count


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
