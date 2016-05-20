package rest

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

const (
	// ContentType is the key for the Content-Type header
	ContentType = "Content-Type"
	// ApplicationJSON is the value for the Content-Type header for JSON
	ApplicationJSON = "application/json"
)

var (
	// JSContentTypeHeader is the http.Header for the "application/json" content type
	JSContentTypeHeader = http.Header(map[string][]string{ContentType: []string{ApplicationJSON}})
)

// Client is a REST client that makes requests to a given path. Implementations will prefix the path with a base URL that is specified when the client implementation is created
type Client interface {
	Do(method string, headers http.Header, body io.Reader, path ...string) (*http.Response, error)
}

type realClient struct {
	baseURL string
	client  *http.Client
}

// NewRealTLSClient creates a new Client that uses a TLS connection to make requests to baseURL
func NewRealTLSClient(baseURL string) Client {
	return &realClient{baseURL: baseURL, client: getTLSClient()}
}

func (r realClient) Do(method string, headers http.Header, body io.Reader, path ...string) (*http.Response, error) {
	strSlice := []string{r.baseURL}
	strSlice = append(strSlice, path...)
	urlStr := strings.Join(strSlice, "/")
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	return r.client.Do(req)
}

func getTLSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	return &http.Client{Transport: tr}
}
