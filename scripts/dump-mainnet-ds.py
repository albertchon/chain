import os
import requests

try:
    os.mkdir("files")
except:
    pass

DS_URL = "http://api-gm-lb.bandchain.org/oracle/data_sources/"
OS_URL = "http://api-gm-lb.bandchain.org/oracle/oracle_scripts/"
DATA_URL = "http://api-gm-lb.bandchain.org/oracle/data/"


def download_data(hash):
    data_url = DATA_URL + hash
    r = requests.get(data_url)
    r.raise_for_status()

    with open("files/" + hash, "wb") as fd:
        for chunk in r.iter_content(chunk_size=128):
            fd.write(chunk)


def prepare_ds(idx, download=False):
    ds_url = DS_URL + str(idx)
    r = requests.get(ds_url)
    r.raise_for_status()

    data = r.json()["result"]
    name = data["name"]
    description = data["description"] if "description" in data else ""
    file = data["filename"]
    if download:
        download_data(file)

    return f"""bandd add-data-source "{name}" "{description}" band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs files/{file}"""


def prepare_os(idx, download=False):
    os_url = OS_URL + str(idx)
    r = requests.get(os_url)
    r.raise_for_status()

    data = r.json()["result"]
    name = data["name"]
    description = data["description"] if "description" in data else ""
    schema = data["schema"]
    file = data["filename"]
    source_code_url = data["source_code_url"] if "source_code_url" in data else ""
    if download:
        download_data(file)

    return f"""bandd add-oracle-script "{name}" "{description}" "{schema}" "{source_code_url}" band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs files/{file}"""


def prepare_mainnet_ds():
    ds = []
    for i in range(1, 100):
        try:
            ds.append(prepare_ds(i, True))
        except Exception as e:
            break

    return ds


def prepare_mainnet_os():
    os = []
    for i in range(1, 100):
        try:
            os.append(prepare_os(i, True))
        except Exception as e:
            break

    return os


if __name__ == "__main__":
    migrate = "#!/bin/bash\n"
    migrate += "\n# Add data sources\n"
    migrate += "\n".join(prepare_mainnet_ds())

    migrate += "\n# Add oracle scripts\n"
    migrate += "\n".join(prepare_mainnet_os())

    s = "scripts/create-mainnet-os-ds.sh"
    os.remove(s)
    f = open(s, "w")

    f.write(migrate)
    f.close()
    os.chmod(s, 436)

    # print(migrate)
