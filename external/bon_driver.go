package external

import (
	"fmt"
	"golang.org/x/exp/maps"
	"os"
	"slices"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

type BonDriverChannelFile struct {
	Headers  []string
	Channels map[MirakurunChannelType][]*BonDriverChannel
}

func NewBonDriverChannelFile() *BonDriverChannelFile {
	return &BonDriverChannelFile{
		Headers: []string{
			"; TVTest チャンネル設定ファイル",
			"; 名称,チューニング空間,チャンネル,リモコン番号,サービスタイプ,サービスID,ネットワークID,TSID,状態",
		},
		Channels: make(map[MirakurunChannelType][]*BonDriverChannel),
	}
}

type BonDriverChannel struct {
	Name               string // 名称
	ChannelIndex       int
	RemoteControlKeyID *uint16 // リモコン番号
	ServiceType        uint8   // サービスタイプ
	ServiceID          uint16  // サービス ID
	NetworkID          uint16  // ネットワーク ID
	TransportStreamID  uint16  // Transport Stream ID
	Enabled            bool    // 状態
}

func (f *BonDriverChannelFile) AddChannel(channelType MirakurunChannelType, channel *BonDriverChannel) {
	f.Channels[channelType] = append(f.Channels[channelType], channel)
}

func (f *BonDriverChannelFile) Write(path string) error {
	var buffer []string
	buffer = append(buffer, f.Headers...)

	var channelTypes []MirakurunChannelType
	for _, channelType := range []MirakurunChannelType{
		MirakurunChannelTypeGR,
		MirakurunChannelTypeBS,
		MirakurunChannelTypeCS,
		MirakurunChannelTypeSKY,
	} {
		if slices.Contains(maps.Keys(f.Channels), channelType) {
			channelTypes = append(channelTypes, channelType)
		}
	}

	var additionalChannelTypes []MirakurunChannelType
	for channelType := range f.Channels {
		if !slices.Contains(channelTypes, channelType) {
			additionalChannelTypes = append(additionalChannelTypes, channelType)
		}
	}

	slices.Sort(additionalChannelTypes)
	channelTypes = append(channelTypes, additionalChannelTypes...)

	for tuningSpaceIndex, channelType := range channelTypes {
		buffer = append(buffer, fmt.Sprintf(";#SPACE(%d,%s)", tuningSpaceIndex, channelType))

		for i, channel := range f.Channels[channelType] {
			var enabled int
			if channel.Enabled {
				enabled = 1
			}

			channelIndex := i
			if channel.ChannelIndex != 0 {
				channelIndex = channel.ChannelIndex
			}

			remoteControlKeyID := channel.ServiceID
			if channel.RemoteControlKeyID != nil && *channel.RemoteControlKeyID > 0 {
				remoteControlKeyID = *channel.RemoteControlKeyID
			}

			buffer = append(buffer,
				fmt.Sprintf(
					"%s,%d,%d,%d,%d,%d,%d,%d,%d",
					channel.Name,
					tuningSpaceIndex,
					channelIndex,
					remoteControlKeyID,
					channel.ServiceType,
					channel.ServiceID,
					channel.NetworkID,
					channel.TransportStreamID,
					enabled,
				),
			)
		}
	}

	encoder := japanese.ShiftJIS.NewEncoder()
	utf8String := strings.Join(buffer, "\r\n")
	shiftJisString, err := encoder.String(utf8String)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(shiftJisString), 0644)
}
