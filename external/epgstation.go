package external

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EPGStationChannel struct {
	ID                 uint64               `json:"id"`
	ServiceID          uint16               `json:"serviceId"`
	NetworkID          uint16               `json:"networkId"`
	Name               string               `json:"name"`
	HalfWidthName      string               `json:"halfWidthName"`
	HasLogoData        bool                 `json:"hasLogoData"`
	ChannelType        MirakurunChannelType `json:"channelType"`
	Channel            string               `json:"channel"`
	Type               *uint8               `json:"type"`
	RemoteControlKeyID *uint16              `json:"remoteControlKeyId"`
}

type EPGStationSchedule struct {
	Channel struct {
		Id          uint64 `json:"id"`
		ServiceId   uint16 `json:"serviceId"`
		NetworkId   uint16 `json:"networkId"`
		Name        string `json:"name"`
		HasLogoData bool   `json:"hasLogoData"`
		ChannelType string `json:"channelType"`
		Type        *uint8 `json:"type"`
	} `json:"channel"`
	Programs []struct {
		Id                 uint64 `json:"id"`
		ChannelId          uint64 `json:"channelId"`
		StartAt            int64  `json:"startAt"`
		EndAt              int64  `json:"endAt"`
		IsFree             bool   `json:"isFree"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		Genre1             int    `json:"genre1"`
		SubGenre1          int    `json:"subGenre1"`
		VideoType          string `json:"videoType"`
		VideoResolution    string `json:"videoResolution"`
		VideoStreamContent int    `json:"videoStreamContent"`
		VideoComponentType int    `json:"videoComponentType"`
		AudioSamplingRate  int    `json:"audioSamplingRate"`
		AudioComponentType int    `json:"audioComponentType"`
		Extended           string `json:"extended"`
		Genre2             int    `json:"genre2"`
		SubGenre2          int    `json:"subGenre2"`
		Genre3             int    `json:"genre3"`
		SubGenre3          int    `json:"subGenre3"`
	} `json:"programs"`
}

type EPGStationClient struct {
	client  *http.Client
	baseURL *url.URL
	header  http.Header
}

func NewEPGStationClient(client *http.Client, baseURL *url.URL, header http.Header) *EPGStationClient {
	return &EPGStationClient{
		client:  client,
		baseURL: baseURL,
		header:  header,
	}
}

func (c *EPGStationClient) GetChannels(ctx context.Context) ([]*EPGStationChannel, error) {
	requestURL := c.baseURL.JoinPath("api", "channels").String()
	slog.InfoContext(ctx, "fetching EPGStation channels", slog.String("url", requestURL))

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	for key, values := range c.header {
		for _, value := range values {
			request.Header.Add(strings.TrimSpace(key), strings.TrimSpace(value))
		}
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data []*EPGStationChannel
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *EPGStationClient) GetChannelSchedules(ctx context.Context, channelID uint64) ([]*EPGStationSchedule, error) {
	u := c.baseURL.JoinPath("api", "schedules", strconv.FormatUint(channelID, 10))
	query := u.Query()
	query.Set("startAt", strconv.FormatInt(time.Now().UnixMilli(), 10))
	query.Set("days", "7")
	query.Set("isHalfWidth", "true")
	u.RawQuery = query.Encode()

	requestURL := u.String()
	slog.InfoContext(ctx, "fetching EPGStation channel schedules", slog.Uint64("channel_id", channelID), slog.String("url", requestURL))

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	for key, values := range c.header {
		for _, value := range values {
			request.Header.Add(strings.TrimSpace(key), strings.TrimSpace(value))
		}
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data []*EPGStationSchedule
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
