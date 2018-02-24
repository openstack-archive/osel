package qualys

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

// Response is a Qualys API response. This wraps the standard http.Response returned from Qualys.
type Response struct {
	*http.Response

	Rate
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is form encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)
	buf := new(bytes.Buffer)

	if method == http.MethodPost {
		buf, err = formPostBody(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.SetBasicAuth(c.Credentials.Username, c.Credentials.Password)
	req.Header.Set(headerUserAgent, userAgent)
	return req, nil
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	response.populateRate()

	return &response
}

// populateRate parses the rate related headers and populates the response Rate.
func (r *Response) populateRate() {
	// TODO - deal with the rest of the headers
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		r.Rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		r.Rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if rateLimitWindow := r.Header.Get(headerRateLimitWindow); rateLimitWindow != "" {
		r.Rate.LimitWindow, _ = strconv.Atoi(rateLimitWindow)
	}
	if waitingPeriod := r.Header.Get(headerRateLimitWait); waitingPeriod != "" {
		r.Rate.WaitingPeriod, _ = strconv.Atoi(waitingPeriod)
	}
	if concurrencyLimit := r.Header.Get(headerConcurrencyLimit); concurrencyLimit != "" {
		r.Rate.ConcurrencyLimit, _ = strconv.Atoi(concurrencyLimit)
	}
	if runningConcurrencyLimit := r.Header.Get(headerConcurrencyLimitRunning); runningConcurrencyLimit != "" {
		r.Rate.CurrentConcurrency, _ = strconv.Atoi(runningConcurrencyLimit)
	}
}

// MakeRequest sends an API request and returns the API response. The API response is XML decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred..
func (c *Client) MakeRequest(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)
	c.Rate = response.Rate

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	bodyContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if v != nil {
		err := xml.Unmarshal(bodyContents, v)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	return fmt.Errorf("Response status is: %v", r.StatusCode)
}

func formPostBody(opt interface{}) (*bytes.Buffer, error) {
	vals, err := query.Values(opt)
	if err != nil {
		return nil, err
	}
	return bytes.NewBufferString(vals.Encode()), nil
}
