package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	Host string
}

func (c *HttpClient) MakeRequest(method string, url string, headers map[string]string, bodyParams []byte, v interface{}) error {
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	endpoint := c.Host + url

	request, err := http.NewRequest(method, endpoint, bytes.NewBuffer(bodyParams))

	if err != nil {
		return err
	}

	for name, value := range headers {
		request.Header.Add(name, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error request: status code: %d", response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &v)
}
