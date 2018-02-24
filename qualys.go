package main

/*

qualys - This file includes all of the logic necessary to interact with the
go-qualys library.  This is extrapolated out so that a QualysInterface
interface can be passed to functions.  Doing this allows testing by mock
classes to be created that can be passed to functions.

Since this is a wrapper around the go-qualys library, this does not need
testing.

*/

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"git.openstack.org/openstack/osel/qualys"
)

// QualysActioner is an interface for an QualysActions class.  Having
// this as an interface allows us to pass in a dummy class for testing that
// just returns mocked data.
type QualysActioner interface {
	InitiateScan([]string) (string, error)
	DropIPv6() bool
}

// QualysActions is a class that handles all interactions directly with Qualys.
// See the comment on QualysActioner for rationale.
type QualysActions struct {
	Options QualysOptions
}

// QualysOptions is a class to convey all of the configurable options for the
// QualysActions class.
type QualysOptions struct {
	DropIPv6       bool
	MinRemaining   int
	ProxyURL       *url.URL
	Password       string
	QualysURL      *url.URL
	ScanOptionName string
	UserName       string
}

// InitiateScan is the main method for the QualysActioner class, it
// makes a call to the Qualys API to start a scan and harvests a scan ID, and
// an optional error string if there is a problem contacting Qualys.
func (s *QualysActions) InitiateScan(targetIPAddresses []string) (string, error) {
	var err error

	// create client with proxy so the qualys service can be accessed
	qualysCreds := qualys.Credentials{
		Username: s.Options.UserName,
		Password: s.Options.Password,
	}
	c, err := qualys.NewClient(&http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(s.Options.ProxyURL)}}, &qualysCreds)
	if err != nil {
		return "", err
	}
	c.BaseURL = s.Options.QualysURL

	// create the options
	opts := qualys.LaunchScanOptions{
		ScanTitle:   "osel",
		ScannerName: "External",
		OptionTitle: s.Options.ScanOptionName,
		IP:          targetIPAddresses,
	}

	// launch the request
	launchScanResponse, err := c.LaunchScan(&opts)
	if err != nil {
		return "", err
	}

	// process the request response
	scanID := launchScanResponse.ScanReference
	remainingQualysRequests := launchScanResponse.RateLimitations.Remaining
	allowedQualysRequests := launchScanResponse.RateLimitations.Limit
	if Debug {
		log.Printf("Qualys Rate Limit: %d of %d total requests remaining, concurrency of %d out of %d, %d seconds remaining in limit window and %d seconds until a request can be made again\n",
			remainingQualysRequests, allowedQualysRequests, launchScanResponse.RateLimitations.CurrentConcurrency,
			launchScanResponse.RateLimitations.ConcurrencyLimit, launchScanResponse.RateLimitations.LimitWindow, launchScanResponse.RateLimitations.WaitingPeriod)
	}
	if launchScanResponse.Text != "" {
		err = errors.New(launchScanResponse.Text)
	}
	if remainingQualysRequests <= s.Options.MinRemaining {
		err = fmt.Errorf("halting Qualys processing!  Only %d Qualys calls remain out of a total of %d.  Waiting for %d seconds before resuming", remainingQualysRequests,
			allowedQualysRequests, launchScanResponse.RateLimitations.LimitWindow)
	}
	return scanID, err
}

// DropIPv6 is an accessor method to allow other code to make decisions based on whether this flag is enabled.
func (s *QualysActions) DropIPv6() bool {
	return s.Options.DropIPv6
}
