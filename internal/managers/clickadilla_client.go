package managers

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
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
}

type ClickadillaClient struct {
	Client *fasthttp.Client
	Host   string
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

func NewClickadillaClient(host string) *ClickadillaClient {
	return &ClickadillaClient{
		Client: &fasthttp.Client{},
		Host:   host,
	}
}

func (c *ClickadillaClient) makeRequest(method string, url string, v interface{}) error {
	var request = fasthttp.AcquireRequest()
	request.SetRequestURI(c.Host + url)

	request.Header.SetMethod(method)

	var response = fasthttp.AcquireResponse()
	err := c.Client.Do(request, response)
	if err != nil {
		return err
	}
	fasthttp.ReleaseRequest(request)
	if statusCode := response.StatusCode(); statusCode != http.StatusOK {
		return fmt.Errorf("storage update error. Response status code client: %d", statusCode)
	}
	defer fasthttp.ReleaseResponse(response)
	return json.Unmarshal(response.Body(), &v)
}

func (c *ClickadillaClient) GetFeeds() ([]Feed, error) {
	response := &FeedsResponse{}

	err := c.makeRequest("GET", "api/billing/v1/feeds", response)

	if err != nil {
		return nil, err
	}
	return response.Feeds, err
}

func (c *ClickadillaClient) GetFeedsTargets() ([]FeedTargers, error) {
	response := struct {
		Targets []FeedTargers `json:"data"`
	}{}

	err := c.makeRequest("GET", "api/billing/v1/feeds-targets", &response)

	if err != nil {
		return nil, err
	}
	return response.Targets, err
}

func (c *ClickadillaClient) GetFeedsSupplySidePlatforms() ([]FeedSupplySidePlatforms, error) {
	response := struct {
		SupplySidePlatforms []FeedSupplySidePlatforms `json:"data"`
	}{}

	err := c.makeRequest("GET", "api/billing/v1/feeds-supply-side-platforms", &response)

	if err != nil {
		return nil, err
	}
	return response.SupplySidePlatforms, err
}

func (c *ClickadillaClient) GetFeedsLabels() ([]FeedLabels, error) {
	response := struct {
		Labels []FeedLabels `json:"data"`
	}{}

	err := c.makeRequest("GET", "api/billing/v1/feeds-labels", &response)

	if err != nil {
		return nil, err
	}
	return response.Labels, err
}

func (c *ClickadillaClient) GetFeedsRtbCategories() ([]FeedRtbCategories, error) {
	response := struct {
		RtbCategories []FeedRtbCategories `json:"data"`
	}{}

	err := c.makeRequest("GET", "api/billing/v1/feeds-rtb-categories", &response)

	if err != nil {
		return nil, err
	}
	return response.RtbCategories, err
}

func (c *ClickadillaClient) GetFeedsNetworks() ([]FeedsNetworks, error) {
	response := make([]FeedsNetworks, 0)

	err := c.makeRequest("GET", "api/billing/v1/feeds-networks", &response)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *ClickadillaClient) GetFeedsAccountManagers() ([]FeedsAccountManagers, error) {
	response := make([]FeedsAccountManagers, 0)

	err := c.makeRequest("GET", "api/billing/v1/feeds-manager-accounts", &response)

	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *ClickadillaClient) GetSupplySidePlatforms() ([]SupplySidePlatform, error) {
	response := &SupplySidePlatformsResponse{}

	err := c.makeRequest("GET", "api/billing/v1/supply-side-platforms", response)
	if err != nil {
		return nil, err
	}
	return response.SupplySidePlatforms, err
}

func (c *ClickadillaClient) GetNetworks() ([]Network, error) {
	response := &NetworksResponse{}

	err := c.makeRequest("GET", "api/billing/v1/networks", response)
	if err != nil {
		return nil, err
	}
	return response.Networks, err
}

func (c *ClickadillaClient) GetDiscrepancies(startDate, endDate time.Time) ([]Discrepancies, error) {
	responseClient := &DiscrepanciesResponse{}
	response := make([]Discrepancies, 0)

	url := fmt.Sprintf("api/billing/v1/feeds/discrepancy-statistics?startDate=%s&endDate=%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	err := c.makeRequest("GET", url, responseClient)
	if err != nil {
		return nil, err
	}

	responseClientByte, err := json.Marshal(responseClient.Discrepancies)

	err = json.Unmarshal(responseClientByte, &response)

	return response, err
}
