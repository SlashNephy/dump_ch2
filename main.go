package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/samber/lo"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"

	"github.com/SlashNephy/dump_ch2/external"
	_ "github.com/SlashNephy/dump_ch2/logger"
)

type BonDriverType string

const (
	BonDriverMirakurun  BonDriverType = "BonDriver_Mirakurun"
	BonDriverMirakc     BonDriverType = "BonDriver_mirakc"
	BonDriverEPGStation BonDriverType = "BonDriver_EPGStation"
)

var opts struct {
	OutputPath           string        `short:"o" long:"output" description:"Output path. default: ${BonDriverType}.ch2"`
	BonDriverType        BonDriverType `short:"t" long:"type" description:"BonDriver_Mirakurun or BonDriver_mirakc, BonDriver_EPGStation?" required:"true"`
	MirakurunScheme      string        `long:"mirakurun-scheme" description:"Mirakurun/mirakc scheme." default:"http"`
	MirakurunHost        string        `long:"mirakurun-host" description:"Mirakurun/mirakc address." default:"127.0.0.1"`
	MirakurunPort        uint16        `long:"mirakurun-port" description:"Mirakurun/mirakc port." default:"40772"`
	EPGStationScheme     string        `long:"epgstation-scheme" description:"EPGStation scheme." default:"http"`
	EPGStationHost       string        `long:"epgstation-host" description:"EPGStation address." default:"127.0.0.1"`
	EPGStationPort       uint16        `long:"epgstation-port" description:"EPGStation port." default:"8888"`
	RequestHeaders       http.Header   `short:"H" long:"header" description:"Request headers."`
	NormalizeServiceName bool          `short:"n" long:"normalize" description:"Normalize service name? (Convert full-width characters to half-width)"`
	StripServices        bool          `short:"s" long:"strip" description:"Exclude unnecessary services?"`
}

func main() {
	ctx := context.Background()

	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	slog.InfoContext(ctx, "application started", slog.Any("opts", opts))

	if err := execute(ctx); err != nil {
		panic(err)
	}
}

func execute(ctx context.Context) error {
	file := external.NewBonDriverChannelFile()

	client := &http.Client{}
	mirakurun := external.NewMirakurunClient(client, &url.URL{Scheme: opts.MirakurunScheme, Host: fmt.Sprintf("%s:%d", opts.MirakurunHost, opts.MirakurunPort)}, opts.RequestHeaders)
	epgstation := external.NewEPGStationClient(client, &url.URL{Scheme: opts.EPGStationScheme, Host: fmt.Sprintf("%s:%d", opts.EPGStationHost, opts.EPGStationPort)}, opts.RequestHeaders)

	services, err := mirakurun.GetServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch services: %w", err)
	}

	transportStreamIDs := make(map[uint64]uint16)
	for _, service := range services {
		if service.TransportStreamID != nil && *service.TransportStreamID != 0 {
			transportStreamIDs[service.ID] = *service.TransportStreamID
		}
	}

	switch opts.BonDriverType {
	case BonDriverMirakurun, BonDriverMirakc:
		serviceGroups := lo.GroupBy(services, func(service *external.MirakurunService) external.MirakurunChannelType {
			return service.Channel.Type
		})

		for channelType, groupedServices := range serviceGroups {
			initialIndices := make(map[uint64]int)
			for i, service := range groupedServices {
				initialIndices[service.ID] = i
			}

			switch channelType {
			case external.MirakurunChannelTypeBS, external.MirakurunChannelTypeCS, external.MirakurunChannelTypeSKY:
				slices.SortStableFunc(groupedServices, func(a, b *external.MirakurunService) int {
					return cmp.Compare(a.ServiceID, b.ServiceID)
				})
			}

			for _, service := range groupedServices {
				serviceType, valid := checkServiceType(service.Name, service.Type)
				enabled := true
				if !valid {
					enabled = false
				} else if opts.StripServices {
					schedules, err := epgstation.GetChannelSchedules(ctx, service.ID)
					if err != nil {
						return fmt.Errorf("failed to fetch schedules: %w", err)
					}

					enabled = slices.ContainsFunc(schedules, func(schedule *external.EPGStationSchedule) bool { return len(schedule.Programs) > 0 })
				}

				name := service.Name
				if opts.NormalizeServiceName {
					name = width.Fold.String(norm.NFKC.String(service.Name))
				}

				transportStreamID, _ := transportStreamIDs[service.ID] //nolint:gosimple
				file.AddChannel(channelType, &external.BonDriverChannel{
					Name:               name,
					ChannelIndex:       initialIndices[service.ID],
					RemoteControlKeyID: service.RemoteControlKeyID,
					ServiceType:        serviceType,
					ServiceID:          service.ServiceID,
					NetworkID:          service.NetworkID,
					TransportStreamID:  transportStreamID,
					Enabled:            enabled,
				})
			}
		}
	case BonDriverEPGStation:
		channels, err := epgstation.GetChannels(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch channels: %w", err)
		}

		channelGroups := lo.GroupBy(channels, func(channel *external.EPGStationChannel) external.MirakurunChannelType {
			return channel.ChannelType
		})

		for channelType, groupedChannels := range channelGroups {
			for _, channel := range groupedChannels {
				serviceType, valid := checkServiceType(channel.Name, channel.Type)
				enabled := true
				if !valid {
					enabled = false
				} else if opts.StripServices {
					schedules, err := epgstation.GetChannelSchedules(ctx, channel.ID)
					if err != nil {
						return fmt.Errorf("failed to fetch schedules: %w", err)
					}

					enabled = slices.ContainsFunc(schedules, func(schedule *external.EPGStationSchedule) bool { return len(schedule.Programs) > 0 })
				}

				name := channel.Name
				if opts.NormalizeServiceName {
					name = norm.NFKC.String(channel.HalfWidthName)
				}

				transportStreamID, _ := transportStreamIDs[channel.ID] //nolint:gosimple
				file.AddChannel(channelType, &external.BonDriverChannel{
					Name:               name,
					RemoteControlKeyID: channel.RemoteControlKeyID,
					ServiceType:        serviceType,
					ServiceID:          channel.ServiceID,
					NetworkID:          channel.NetworkID,
					TransportStreamID:  transportStreamID,
					Enabled:            enabled,
				})
			}
		}
	default:
		return fmt.Errorf("invalid BonDriverType: %s", opts.BonDriverType)
	}

	path := opts.OutputPath
	if path == "" {
		path = fmt.Sprintf("%s.ch2", opts.BonDriverType)
	}

	if err = file.Write(path); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	slog.InfoContext(ctx, "wrote file", slog.String("path", path))
	return nil
}

func checkServiceType(serviceName string, serviceType *uint8) (uint8, bool) {
	if serviceType == nil {
		if strings.Contains(serviceName, "ラジオ") {
			return 0x02, true // 音声サービス
		}
		return 0x01, true // 映像サービス
	}

	return *serviceType, isVideoOrAudioService(*serviceType)
}

func isVideoOrAudioService(serviceType uint8) bool {
	switch serviceType {
	case 0x01, 0x02, 0xa1, 0xa2, 0xa5, 0xa6, 0xad:
		return true
	default:
		return false
	}
}
