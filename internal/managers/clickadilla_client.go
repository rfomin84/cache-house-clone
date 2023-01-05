package managers

import (
	"encoding/json"
	"fmt"
	"github.com/clickadilla/cache-house/pkg/httpclient"
	"net/http"
	"time"
)

type ClickadillaClientInterface interface {
	GetFeeds() ([]Feed, error)
	GetFeedsNetworks() ([]FeedsNetworks, error)
	GetFeedsAccountManagers() ([]FeedsAccountManagers, error)
	GetFeedsTargets() ([]FeedTargers, error)
	GetFeedsSupplySidePlatforms() ([]FeedSupplySidePlatforms, error)
	GetSupplySidePlatforms() ([]SupplySidePlatform, error)
	GetFeedsLabels() ([]FeedLabels, error)
	GetFeedsRtbCategories() ([]FeedRtbCategories, error)
	GetNetworks() ([]Network, error)
	GetDiscrepancies(startDate, endDate time.Time) ([]Discrepancies, error)
	GetAllFeeds() ([]AllFeeds, error)
}

type ClickadillaClient struct {
	Client   *httpclient.HttpClient
	apiToken string
}

func NewClickadillaClient(host, token string) *ClickadillaClient {
	return &ClickadillaClient{
		Client: &httpclient.HttpClient{
			Host: host,
		},
		apiToken: token,
	}
}

type FeedsResponse struct {
	Feeds []Feed `json:"data"`
}

type SupplySidePlatformsResponse struct {
	SupplySidePlatforms []SupplySidePlatform `json:"data"`
}

type NetworksResponse struct {
	Networks []Network `json:"data"`
}

type DiscrepanciesResponse struct {
	Discrepancies []struct {
		Date                 string  `json:"date"`
		FeedId               int     `json:"feed_id"`
		Discrepancy          float64 `json:"discrepancy"`
		IsDemandSidePlatform bool    `json:"is_demand_side_platform"`
	} `json:"data"`
}

func (c *ClickadillaClient) GetFeeds() ([]Feed, error) {
	response := &FeedsResponse{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds", map[string]string{}, nil, response)

	if err != nil {
		return nil, err
	}
	return response.Feeds, err
}

func (c *ClickadillaClient) GetFeedsTargets() ([]FeedTargers, error) {
	response := struct {
		Targets []FeedTargers `json:"data"`
	}{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-targets", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response.Targets, err
}

func (c *ClickadillaClient) GetFeedsSupplySidePlatforms() ([]FeedSupplySidePlatforms, error) {
	response := struct {
		SupplySidePlatforms []FeedSupplySidePlatforms `json:"data"`
	}{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-supply-side-platforms", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response.SupplySidePlatforms, err
}

func (c *ClickadillaClient) GetFeedsLabels() ([]FeedLabels, error) {
	response := struct {
		Labels []FeedLabels `json:"data"`
	}{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-labels", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response.Labels, err
}

func (c *ClickadillaClient) GetFeedsRtbCategories() ([]FeedRtbCategories, error) {
	response := struct {
		RtbCategories []FeedRtbCategories `json:"data"`
	}{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-rtb-categories", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response.RtbCategories, err
}

func (c *ClickadillaClient) GetFeedsNetworks() ([]FeedsNetworks, error) {
	response := make([]FeedsNetworks, 0)

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-networks", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *ClickadillaClient) GetFeedsAccountManagers() ([]FeedsAccountManagers, error) {
	response := make([]FeedsAccountManagers, 0)

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/feeds-manager-accounts", map[string]string{}, nil, &response)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *ClickadillaClient) GetSupplySidePlatforms() ([]SupplySidePlatform, error) {
	response := &SupplySidePlatformsResponse{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/supply-side-platforms", map[string]string{}, nil, response)
	if err != nil {
		return nil, err
	}
	return response.SupplySidePlatforms, err
}

func (c *ClickadillaClient) GetNetworks() ([]Network, error) {
	response := &NetworksResponse{}

	err := c.Client.MakeRequest(http.MethodGet, "api/billing/v1/networks", map[string]string{}, nil, response)
	if err != nil {
		return nil, err
	}
	return response.Networks, err
}

func (c *ClickadillaClient) GetDiscrepancies(startDate, endDate time.Time) ([]Discrepancies, error) {
	responseClient := &DiscrepanciesResponse{}
	response := make([]Discrepancies, 0)

	url := fmt.Sprintf("api/billing/v1/feeds/discrepancy-statistics?startDate=%s&endDate=%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	err := c.Client.MakeRequest(http.MethodGet, url, map[string]string{}, nil, responseClient)
	if err != nil {
		return nil, err
	}

	responseClientByte, err := json.Marshal(responseClient.Discrepancies)

	err = json.Unmarshal(responseClientByte, &response)

	return response, err
}

func (c *ClickadillaClient) GetAllFeeds() ([]AllFeeds, error) {
	response := make([]AllFeeds, 0)

	url := fmt.Sprintf("api/internal/v1/feeds/list-for-gather-statistics")
	headers := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.apiToken,
	}
	err := c.Client.MakeRequest(http.MethodPost, url, headers, nil, &response)

	if err != nil {
		return nil, err
	}

	return response, err
}
