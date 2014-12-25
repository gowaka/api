package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/httgo/interfaces/httpclient"
	"net/http"
	"net/url"
)

var (
	ErrBlankAPIKey = errors.New("api key: cannot be blank")
)

type Client struct {
	Client httpclient.Interface
	ApiKey string
}

func NewClient(key string, opts ...func(*Client)) (*Client, error) {
	if key == "" {
		return nil, ErrBlankAPIKey
	}

	c := &Client{
		Client: http.DefaultClient,
		ApiKey: key,
	}
	for _, v := range opts {
		v(c)
	}

	return c, nil
}

func enc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func authorize(req *http.Request, c *Client) {
	req.Header.Add("Authorization", "Basic "+enc(c.ApiKey))
}

// Do calls Client.Do with authorization headers
func (c Client) Do(req *http.Request) (*http.Response, error) {
	authorize(req, &c)
	return c.Client.Do(req)
}

const (
	SCHEME = "https"
	HOST   = "wakatime.com"
)

func URL(path string) *url.URL {
	return &url.URL{
		Scheme: SCHEME,
		Host:   HOST,
		Path:   path,
	}
}

// Get calls a resource and decodes the response body to d
func (c Client) Get(r Resource, d interface{}) error {
	resp, err := c.GetResp(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(d)
	if err != nil {
		return err
	}

	return nil
}

// GetResp calls a resource and returns the raw response
func (c Client) GetResp(r Resource) (*http.Response, error) {
	req, err := NewRequest(r)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// NewRequest prepares a request from a resource
func NewRequest(r Resource) (*http.Request, error) {
	u := URL(r.Path())
	if q := r.Filter(); q != nil {
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
