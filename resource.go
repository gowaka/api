package api

// Resource represents the API resource interface
type Resource interface {
	Get(interface{}) error
	Path() string
}
