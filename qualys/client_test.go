package qualys

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux *http.ServeMux

	client *Client

	server *httptest.Server

	creds *Credentials
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	creds = &Credentials{Username: "bogus", Password: "bogus"}
	client, _ = NewClient(nil, creds)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parseForm(): %v", err)
	}

	if !reflect.DeepEqual(expected, r.Form) {
		t.Errorf("Request parameters = %v, expected %v", r.Form, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	if c.BaseURL == nil || c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, defaultBaseURL)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if c.UserAgent != userAgent {
		t.Errorf("NewClick UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultBaseURL(t, c)
	testClientDefaultUserAgent(t, c)
}

func TestNewClient(t *testing.T) {
	c, _ := NewClient(nil, creds)
	testClientDefaults(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(nil, creds)

	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	testClientDefaults(t, c)
}

func TestNewClientWithoutCredentials(t *testing.T) {
	_, err := NewClient(nil, nil)

	if err == nil {
		t.Errorf("NewClient() expected error when Credentials are not set")
	}
}

func TestCustomUserAgent(t *testing.T) {
	c, err := New(nil, creds, SetUserAgent("testing"))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s+%s", "testing", userAgent)
	if got := c.UserAgent; got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestAddURLParameters(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected string
		opts     *ListAssetGroupOptions
		isErr    bool
	}{
		{
			name:     "addURLParameters",
			path:     "/asset/group/",
			expected: "/asset/group/?ids=1",
			opts:     &ListAssetGroupOptions{Ids: []string{"1"}},
			isErr:    false,
		},
		{
			name:     "addURLParameters with slice parameter",
			path:     "/asset/group/",
			expected: "/asset/group/?ids=1,2",
			opts:     &ListAssetGroupOptions{Ids: []string{"1", "2"}},
			isErr:    false,
		},
	}

	for _, c := range cases {
		got, err := addURLParameters(c.path, c.opts)
		if c.isErr && err == nil {
			t.Errorf("%q expected error but none was encountered", c.name)
			continue
		}

		if !c.isErr && err != nil {
			t.Errorf("%q unexpected error: %v", c.name, err)
			continue
		}

		gotURL, err := url.Parse(got)
		if err != nil {
			t.Errorf("%q unable to parse returned URL", c.name)
			continue
		}

		expectedURL, err := url.Parse(c.expected)
		if err != nil {
			t.Errorf("%q unable to parse expected URL", c.name)
			continue
		}

		if g, e := gotURL.Path, expectedURL.Path; g != e {
			t.Errorf("%q path = %q; expected %q", c.name, g, e)
			continue
		}

		if g, e := gotURL.Query(), expectedURL.Query(); !reflect.DeepEqual(g, e) {
			t.Errorf("%q query = %#v; expected %#v", c.name, g, e)
			continue
		}
	}
}

func TestFormPostBody(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		opts     *AddIPsToGroupOptions
	}{
		{
			name:     "formPostBody single IP",
			expected: "add_ips=10.10.10.10&id=1234",
			opts:     &AddIPsToGroupOptions{GroupID: "1234", IPs: []string{"10.10.10.10"}},
		},
		{
			name:     "formPostBody multi IP",
			expected: "add_ips=10.10.10.10%2C10.10.10.11&id=1234",
			opts:     &AddIPsToGroupOptions{GroupID: "1234", IPs: []string{"10.10.10.10", "10.10.10.11"}},
		},
	}
	for _, c := range cases {
		got, _ := formPostBody(c.opts)
		buf := new(bytes.Buffer)
		buf.ReadFrom(got)
		bodyString := buf.String()
		if c.expected != bodyString {
			t.Errorf("%q expected %s but got %s", c.name, c.expected, bodyString)
		}
	}
}
