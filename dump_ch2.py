import argparse
import json
import itertools
import sys
import urllib.request
import unicodedata

def call_json_api(url):
    with urllib.request.urlopen(url) as response:
        content = response.read()

        return json.loads(content.decode())

def get_service_type(service):
    if "ラジオ" in service["name"]:
        return 2  # 音声サービス

    return 1  # 映像サービス

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-a", "--address", help="specify Mirakurun/mirakc address.", default="127.0.0.1")
    parser.add_argument("-p", "--port", help="specify Mirakurun/mirakc port.", type=int, default=40772)
    parser.add_argument("-o", "--output", help="specify Output file path.", default="BonDriver_EPGStation.ch2")
    parser.add_argument("-n", "--normalize", help="normalize service name? (convert full-width chars to half-width)", action="store_true")
    args = parser.parse_args()

    url = f"http://{args.address}:{args.port}/api/services"
    print(f"GET {url}")
    response = call_json_api(url)

    output = []

    output.append("; TVTest チャンネル設定ファイル")
    output.append("; 名称,チューニング空間,チャンネル,リモコン番号,サービスタイプ,サービスID,ネットワークID,TSID,状態")

    channel_types_order = {"GR": 0, "BS": 1, "CS": 2, "SKY": 3}
    response.sort(key=lambda x: channel_types_order[x["channel"]["type"]])
    grouped_services = itertools.groupby(response, lambda x: x["channel"]["type"])

    for i, x in enumerate(grouped_services):
        channelType, services = x
        if channelType not in ["BS", "CS"]:
            services = sorted(services, key=lambda x: x["remoteControlKeyId"])
        else:
            services = sorted(services, key=lambda x: x["serviceId"])

        output.append(f";#SPACE({i},{channelType})")

        for j, service in enumerate(services):
            output.append(f"{unicodedata.normalize('NFKC', service['name']) if args.normalize else service['name']},{i},{j},{service['remoteControlKeyId'] if service['remoteControlKeyId'] > 0 else service['serviceId']},{get_service_type(service)},{service['serviceId']},{service['networkId']},{service['transportStreamId']},1")

    with open(args.output, "w", encoding="cp932") as f:
        f.write("\r\n".join(output))
