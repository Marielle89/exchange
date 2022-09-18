package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client for the ExchangeRate-API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
}

const (
	schema  = "https"
	host    = "v6.exchangerate-api.com"
	version = "v6"
	timeout = 3
)

// NewClient creates a Client
func NewClient(options ...func(*Client)) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: timeout * time.Second,
		},
		baseURL: &url.URL{
			Scheme: schema,
			Host:   host,
			Path:   version,
		},
		apiKey: os.Getenv("EXCHANGE_RATE_API_KEY"),
	}

	for _, do := range options {
		do(client)
	}

	return client
}

// ApiKey sets the api key used by connect to the ExchangeRate-API
func ApiKey(ak string) func(*Client) {
	return func(c *Client) {
		c.apiKey = ak
	}
}

// Rate returns the exchange rate from base currency to target currency
func (client *Client) Rate(ctx context.Context, base Currency, target Currency) (Rate, error) {
	response, err := client.get(ctx, fmt.Sprintf("/pair/%s/%s", base, target))
	if err != nil {
		return 0, err
	}

	return response.Rate, nil
}

// Amount returns the exchange rate from base currency to target currency as a conversion of the amount you supplied
func (client *Client) Amount(ctx context.Context, base Currency, target Currency, amount Amount) (Amount, error) {
	response, err := client.get(ctx, fmt.Sprintf("/pair/%s/%s/%f", base, target, amount))

	if err != nil {
		return 0, err
	}

	return response.Amount, nil
}

func (client *Client) get(ctx context.Context, path string) (*Response, error) {
	request, requestErr := client.request(ctx, path)
	if requestErr != nil {
		return nil, requestErr
	}

	response, responseErr := client.do(request)
	if requestErr != nil {
		return nil, responseErr
	}

	return response, nil
}

func (client *Client) do(request *http.Request) (*Response, error) {
	resp, requestErr := client.httpClient.Do(request)
	if requestErr != nil {
		return nil, requestErr
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 64)
		_ = resp.Body.Close()
	}()

	if responseErr := responseError(resp); responseErr != nil {
		return nil, responseErr
	}

	var response Response

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if err := responseApiError(response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (client *Client) request(ctx context.Context, path string) (*http.Request, error) {
	ref, err := url.Parse(client.baseURL.Path + path)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET", client.baseURL.ResolveReference(ref).String(), nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	request.Header.Add("Authorization", "Bearer "+client.apiKey)

	return request, nil
}
