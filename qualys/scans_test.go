package qualys

import (
	"net/http"
	"net/url"
	"testing"
	"time"
)

var devCreds Credentials = Credentials{
	Username: "cmcas_ae2",
	Password: "D02debLYko",
}

func TestLiveScan(t *testing.T) {
	// create client
	c, clientErr := NewClient(&http.Client{}, &devCreds)

	if clientErr != nil {
		t.Error(clientErr)
	}

	if baseURL, err := url.Parse("https://qualysapi.qualys.com/api/2.0/fo/scan/"); err != nil {
		t.Error(err)
	} else {
		c.BaseURL = baseURL
	}

	// create the options
	opts := LaunchScanOptions{
		ScanTitle:   "hello_world",
		ScannerName: "External",
		// OptionID:    923922,
		OptionTitle: "Elastic Cloud Option Profile with Password Guessing",
		IP:          []string{"96.119.99.178"},
	}

	// launch the request

	launchScanResponse, err := c.LaunchScan(&opts)

	if err != nil {
		t.Error(err)
	}

	// not sure if necessary
	time.Sleep(time.Minute * 1)

	//time to poll the scan results
	pollOpts := PollScanOptions{
		ScanRef: launchScanResponse.ScanReference,
	}

	_, pollRespErr := c.PollScanResults(&pollOpts)

	if pollRespErr != nil {
		t.Error(pollRespErr)
	}

	// now need to keep polling until the results are all in...

	resultsOptions := CompletedScanOptions{
		ScanRef: launchScanResponse.ScanReference,
	}

	_, resultsRespErr := c.GetScanResults(&resultsOptions)

	if resultsRespErr != nil {
		t.Error(resultsRespErr)
	}
}
