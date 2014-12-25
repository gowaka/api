package api

import (
	"github.com/httgo/mock"
	"gopkg.in/nowk/assert.v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLReturnsAWakaTimeURL(t *testing.T) {
	u := URL("/path/to/resource")
	assert.Equal(t, "https://wakatime.com/path/to/resource", u.String())
}

func TestAPIAuthorizesRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {})
	mo := &mock.Mock{
		Testing: t,
		Ts:      httptest.NewUnstartedServer(mux),
	}
	mo.Start()
	defer mo.Done()

	api, err := NewClient("12345", func(c *Client) {
		c.Client = mo
	})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", mo.Ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := api.Do(req); err != nil {
		t.Fatal(err)
	}

	reqs := mo.History("GET", mo.Ts.URL)
	assert.Equal(t, 1, len(reqs))

	r := reqs[0]
	assert.Equal(t, "Basic MTIzNDU=", r.Header.Get("Authorization"))
}

type resource struct{}

func (r resource) Path() string {
	return "/foo/bar"
}

func TestBuildRequest(t *testing.T) {
	req, err := NewRequest(&resource{})
	if err != nil {
		t.Fatal(err)
	}

	u := req.URL
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "wakatime.com", u.Host)
	assert.Equal(t, "/foo/bar", u.Path)
}
