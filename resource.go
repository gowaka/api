package api

import (
	"net/url"
)

// Resource represents the API resource interface
type Resource interface {
	Path() string
	Filter() *url.Values
}
