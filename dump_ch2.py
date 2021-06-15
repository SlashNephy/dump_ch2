import argparse
import itertools
import json
import time
import urllib.request


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
    parser.add_argument("-ma", "--mirakurun-address", help="specify Mirakurun/mirakc address.", default="127.0.0.1")
    parser.add_argument("-mp", "--mirakurun-port", help="specify Mirakurun/mirakc port.", type=int, default=40772)
    parser.add_argument("-ea", "--epgstation-address", help="specify EPGStation address.", default="127.0.0.1")
    parser.add_argument("-ep", "--epgstation-port", help="specify EPGStation port.", type=int, default=8888)
    parser.add_argument("-o", "--output", help="specify output file path.", default="BonDriver_EPGStation.ch2")
    parser.add_argument("-n", "--normalize", help="normalize service name? (convert full-width chars to half-width)", action="store_true")
    parser.add_argument("-s", "--strip", help="exclude unnecessary services?", action="store_true")
    args = parser.parse_args()

    mirakurun_services_url = f"http://{args.mirakurun_address.strip()}:{args.mirakurun_port}/api/services"
    print(f"GET {mirakurun_services_url}")
    mirakurun_transport_stream_ids = {
        x["id"]: x["transportStreamId"]
        for x in call_json_api(mirakurun_services_url)
    }

    epgstation_channels_url = f"http://{args.epgstation_address.strip()}:{args.epgstation_port}/api/channels"
    print(f"GET {epgstation_channels_url}")
    epgstation_channels_response = call_json_api(epgstation_channels_url)

    output = []

    output.append("; TVTest チャンネル設定ファイル")
    output.append("; 名称,チューニング空間,チャンネル,リモコン番号,サービスタイプ,サービスID,ネットワークID,TSID,状態")

    grouped_channels = itertools.groupby(epgstation_channels_response, lambda x: x["channelType"])
    for i, x in enumerate(grouped_channels):
        channelType, channels = x

        output.append(f";#SPACE({i},{channelType})")

        for j, channel in enumerate(channels):
            enabled = True
            if args.strip:
                epgstation_schedules_url = f"http://{args.epgstation_address.strip()}:{args.epgstation_port}/api/schedules/{channel['id']}?startAt={int(time.time() * 1000)}&days=7&isHalfWidth=true"
                print(f"GET {epgstation_schedules_url}")
                epgstation_schedules_response = call_json_api(epgstation_schedules_url)

                if all([len(x["programs"]) == 0 for x in epgstation_schedules_response]):
                    enabled = False

            # TSID は EPGStation API から取得できないので Mirakurun のデータを使う
            transport_stream_id = mirakurun_transport_stream_ids.get(channel["id"]) or 0

            output.append(f"{channel['halfWidthName'] if args.normalize else channel['name']},{i},{j},{channel['remoteControlKeyId'] if channel['remoteControlKeyId'] > 0 else channel['serviceId']},{get_service_type(channel)},{channel['serviceId']},{channel['networkId']},{transport_stream_id},{int(enabled)}")

    with open(args.output, "w", encoding="cp932") as f:
        f.write("\r\n".join(output))
    print(f"Output to {args.output}")
