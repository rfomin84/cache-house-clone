package managers

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ClickadillaClientInterface interface {
	GetFeeds() ([]Feed, error)
	GetSupplySidePlatforms() ([]SupplySidePlatform, error)
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

func (c *ClickadillaClient) GetSupplySidePlatforms() ([]SupplySidePlatform, error) {
	response := &SupplySidePlatformsResponse{}

	err := c.makeRequest("GET", "api/billing/v1/supply-side-platforms", response)
	if err != nil {
		return nil, err
	}
	return response.SupplySidePlatforms, err
}
