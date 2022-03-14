package managers

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type ClickadillaClientInterface interface {
	GetFeeds() ([]Feed, error)
}

type ClickadillaClient struct {
	Client   *fasthttp.Client
	Host     string
	ApiToken string
}

type feedsResponse struct {
	Feeds []Feed `json:"data"`
}

func NewClickadillaClient(host string, apiToken string) *ClickadillaClient {
	return &ClickadillaClient{
		Client:   &fasthttp.Client{},
		Host:     host,
		ApiToken: apiToken,
	}
}

func (c *ClickadillaClient) makeRequest(method string, url string, v interface{}) error {
	var request = fasthttp.AcquireRequest()
	request.SetRequestURI(c.Host + url)

	request.Header.SetMethod(method)
	request.Header.Add("Authorization", "Bearer "+c.ApiToken)

	var response = fasthttp.AcquireResponse()
	err := c.Client.Do(request, response)
	if err != nil {
		return err
	}
	fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	return json.Unmarshal(response.Body(), &v)
}

func (c *ClickadillaClient) GetFeeds() ([]Feed, error) {
	response := &feedsResponse{}

	err := c.makeRequest("GET", "api/billing/v1/feeds", response)

	if err != nil {
		return nil, err
	}
	return response.Feeds, err
}
