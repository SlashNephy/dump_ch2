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
usage: dump_ch2.py [-h] [-a ADDRESS] [-p PORT] [-o OUTPUT] [-n]

optional arguments:
  -h, --help            show this help message and exit
  -a ADDRESS, --address ADDRESS
                        specify Mirakurun/mirakc address.
  -p PORT, --port PORT  specify Mirakurun/mirakc port.
  -o OUTPUT, --output OUTPUT
                        specify Output file path.
  -n, --normalize       normalize service name? (convert full-width chars to half-width)
```
