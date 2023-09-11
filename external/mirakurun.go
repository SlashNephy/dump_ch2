package external

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type MirakurunChannelType string

const (
	MirakurunChannelTypeGR  MirakurunChannelType = "GR"
	MirakurunChannelTypeBS  MirakurunChannelType = "BS"
	MirakurunChannelTypeCS  MirakurunChannelType = "CS"
	MirakurunChannelTypeSKY MirakurunChannelType = "SKY"
)

type MirakurunService struct {
	ID                 uint64  `json:"id"`
	ServiceID          uint16  `json:"serviceId"`
	NetworkID          uint16  `json:"networkId"`
	TransportStreamID  *uint16 `json:"transportStreamId"`
	Name               string  `json:"name"`
	Type               *uint8  `json:"type"`
	LogoId             int     `json:"logoId,omitempty"`
	RemoteControlKeyID *uint16 `json:"remoteControlKeyId"`
	EpgReady           bool    `json:"epgReady"`
	EpgUpdatedAt       int64   `json:"epgUpdatedAt"`
	Channel            struct {
		Type    MirakurunChannelType `json:"type"`
		Channel string               `json:"channel"`
	} `json:"channel"`
	HasLogoData bool    `json:"hasLogoData"`
	LogoData    *string `json:"logoData"`
}

type MirakurunClient struct {
	client  *http.Client
	baseURL *url.URL
	header  http.Header
}

func NewMirakurunClient(client *http.Client, baseURL *url.URL, header http.Header) *MirakurunClient {
	return &MirakurunClient{
		client:  client,
		baseURL: baseURL,
		header:  header,
	}
}

func (c *MirakurunClient) GetServices(ctx context.Context) ([]*MirakurunService, error) {
	requestURL := c.baseURL.JoinPath("api", "services").String()
	slog.InfoContext(ctx, "fetching Mirakurun services", slog.String("url", requestURL))

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

	var data []*MirakurunService
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
