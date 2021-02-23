#!/usr/bin/env python3

import json
import urllib.request
import sys

ANU_URL = "https://qrng.anu.edu.au/API/jsonI.php?length={}&type=uint8"


def make_json_request(url):
    return json.loads(urllib.request.urlopen(url).read())


def main(length):
    random = make_json_request(ANU_URL.format(length))
    return ", ".join(str(x) for x in random["data"])


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
