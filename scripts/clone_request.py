import requests
import sys
import base64
from dataclasses import dataclass
from random import sample


# rid = sys.argv[1]
# memo = sys.argv[2] if len(sys.argv) >= 3 else ""
REQ_URL = "http://api-gm-lb.bandchain.org/oracle/requests/"


@dataclass
class Entry:
    oid: int
    calldata: str
    ask_count: int = 16
    min_count: int = 10
    memo: str = ""
    gas: int = 1590000

    def gen_command(self):
        return f"""sleep 6 && yes | bandcli tx oracle request {self.oid} {self.ask_count} {self.min_count} -c {self.calldata} -m requester --from requester$1 --keyring-backend test --chain-id bandchain --gas {self.gas} --memo '{self.memo}';"""


def get_request(id):
    req_url = REQ_URL + str(id)
    r = requests.get(req_url)
    r.raise_for_status()

    result = r.json()["result"]

    return (
        result["request"]["oracle_script_id"],
        base64.b64decode(result["request"]["calldata"]).hex(),
    )


def get_entry(id):
    oid, calldata = get_request(id)

    return Entry(oid, calldata)


# start end sampling
start = int(sys.argv[1])
end = int(sys.argv[2])
sampling = int(sys.argv[3])
ask_count = int(sys.argv[4])
min_count = int(sys.argv[5])
memo = sys.argv[6]

ids = [i for i in range(start, end)]
ids = sample(ids, sampling)

for id in ids:
    e = get_entry(id)
    e.ask_count = ask_count
    e.min_count = min_count
    e.memo = memo
    print(e.gen_command())
