# dump_ch2
BonDriver_EPGStation 用の TVTest のチャンネル設定ファイルを CLI で生成するスクリプト (要 Mirakurun or mirakc, EPGStation 環境)

## Requirements

- Mirakurun or mirakc
- EPGStation
- Python 3.6 or above
- TVTest w/ BonDriver_EPGStation

## Usage

```console
foo@bar:~$ python dump_ch2.py --help
usage: dump_ch2.py [-h] [-ma MIRAKURUN_ADDRESS]
                   [-mp MIRAKURUN_PORT]
                   [-ea EPGSTATION_ADDRESS]
                   [-ep EPGSTATION_PORT]
                   [-o OUTPUT] [-n]

optional arguments:
  -h, --help            show this help message
                        and exit
  -ma MIRAKURUN_ADDRESS, --mirakurun-address MIRAKURUN_ADDRESS
                        specify Mirakurun/mirakc
                        address.
  -mp MIRAKURUN_PORT, --mirakurun-port MIRAKURUN_PORT
                        specify Mirakurun/mirakc
                        port.
  -ea EPGSTATION_ADDRESS, --epgstation-address EPGSTATION_ADDRESS
                        specify EPGStation
                        address.
  -ep EPGSTATION_PORT, --epgstation-port EPGSTATION_PORT
                        specify EPGStation port.
  -o OUTPUT, --output OUTPUT
                        specify output file path.
  -n, --normalize       normalize service name?
                        (convert full-width chars
                        to half-width)
```
