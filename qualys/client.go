package qualys

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://qualysapi.qualys.com/api/2.0/fo/"
	userAgent      = "go-qualys"
	mediaType      = "application/xml"

	headerUserAgent               = "X-Requested-With"
	headerRateLimit               = "X-RateLimit-Limit"
	headerRateLimitWindow         = "X-RateLimit-Window-Sec"
	headerRateRemaining           = "X-RateLimit-Remaining"
	headerRateLimitWait           = "X-RateLimit-ToWait-Sec"
	headerConcurrencyLimit        = "X-Concurrency-Limit-Limit"
	headerConcurrencyLimitRunning = "X-Concurrency-Limit-Running"
)

// Client for Qualys API
type Client struct {
	// Credentials used to authenticate to the Qualys API
	Credentials *Credentials

	// HTTP client used to communicate with the Qualys API
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call.
	Rate Rate

	// Services used for communicating with the API
	Assets AssetsService
}

// Rate contains the rate limit for the current client.
type Rate struct {
	// The number of requests within the limit window of seconds the client is allowed
	Limit int

	// The number of seconds remaining in the limit window
	LimitWindow int

	// The number of remaining requests the client can make during the limit window period
	Remaining int

	// The number of seconds to wait before requests can be made again -- headerRateLimitWait
	WaitingPeriod int

	// The number of API calls permitted to be executed concurrrently
	ConcurrencyLimit int

	// The number of API calls currently running
	CurrentConcurrency int
}

// Credentials holds the credentials and endpoint for the Qualys Client
type Credentials struct {
	Username string
	Password string
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new API client instance.
func New(httpClient *http.Client, credentials *Credentials, opts ...ClientOpt) (*Client, error) {
	c, err := NewClient(httpClient, credentials)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// NewClient returns a new Qualys API client.
func NewClient(httpClient *http.Client, credentials *Credentials) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if credentials == nil || credentials.Username == "" || credentials.Password == "" {
		return nil, fmt.Errorf("Credentials must be provided")
	}

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, Credentials: credentials, BaseURL: baseURL, UserAgent: userAgent}

	c.Assets = &AssetsServiceOp{client: c}
	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s+%s", ua, c.UserAgent)
		return nil
	}
}
