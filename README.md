# dump_ch2

BonDriver 用の TVTest のチャンネル設定ファイルを CLI で生成するスクリプト

## Requirements

- Mirakurun or mirakc
- (任意) EPGStation
- TVTest
  - いずれかの BonDriver
    - BonDriver_Mirakurun
    - BonDriver_mirakc
    - BonDriver_EPGStation

## Usage

```console
$ go run . --help

Usage:
  dump_ch2 [OPTIONS]

Application Options:
  -o, --output:             Output path. default: ${BonDriverType}.ch2
  -t, --type:               BonDriver_Mirakurun or BonDriver_mirakc, BonDriver_EPGStation?
      --mirakurun-scheme:   Mirakurun/mirakc scheme. (default: http)
      --mirakurun-host:     Mirakurun/mirakc address. (default: 127.0.0.1)
      --mirakurun-port:     Mirakurun/mirakc port. (default: 40772)
      --epgstation-scheme:  EPGStation scheme. (default: http)
      --epgstation-host:    EPGStation address. (default: 127.0.0.1)
      --epgstation-port:    EPGStation port. (default: 8888)
  -H, --header:             Request headers.
  -n, --normalize           Normalize service name? (Convert full-width characters to half-width)
  -s, --strip               Exclude unnecessary services?

Help Options:
  -h, --help                Show this help message
```
