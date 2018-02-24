package qualys

import (
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
)

var ErrMalformedResponse error = errors.New("malformed xml response from server")

// LaunchScanResponse is the expected response for a scan launch
// xml tag: simple_return
// DateTime is the date and time the scan was issued
// It includes a key/value map of pertinent information.
type LaunchScanResponse struct {
	XMLName  xml.Name `xml:"SIMPLE_RETURN"`
	Value    []string `xml:"RESPONSE>ITEM_LIST>ITEM>VALUE"`
	Key      []string `xml:"RESPONSE>ITEM_LIST>ITEM>KEY"`
	Datetime string   `xml:"RESPONSE>DATETIME"`
	Text     string   `xml:"RESPONSE>TEXT"`

	ScanReference string

	RateLimitations Rate
}

type LaunchScanOptions struct {
	ScanTitle   string   `url:"scan_title"`
	OptionID    int64    `url:"option_id,omitempty"`
	OptionTitle string   `url:"option_title"`
	ScannerName string   `url:"iscanner_name"`
	IP          []string `url:"ip"`
	Action      string   `url:"action"`
}

// ScanLaunch will try and run a scan on demand
// has certain parameters that must be filled in
func (client *Client) LaunchScan(options *LaunchScanOptions) (*LaunchScanResponse, error) {
	options.Action = "launch"

	urlString, err := addURLParameters(client.BaseURL.String(), options)

	if err != nil {
		return nil, err
	}

	// i don't think there's any header data here?
	req, err := client.NewRequest(http.MethodPost, urlString, nil)
	if err != nil {
		return nil, err
	}

	var resp LaunchScanResponse

	response, err := client.MakeRequest(req, &resp)

	if err != nil {
		return nil, err
	}

	resp.RateLimitations = response.Rate

	for key, val := range resp.Key {
		if strings.ToUpper(val) == "REFERENCE" {
			// should check len first
			if len(resp.Value) > key {
				resp.ScanReference = resp.Value[key]
			} else {
				return nil, ErrMalformedResponse
			}
		}
	}

	return &resp, nil
}

type PollScanResponse struct {
	XMLName             xml.Name `xml:"SCAN_LIST_OUTPUT"`
	Datetime            string   `xml:"RESPONSE>DATETIME"`
	Processing_priority string   `xml:"RESPONSE>SCAN_LIST>SCAN>PROCESSING_PRIORITY"`
	Processed           string   `xml:"RESPONSE>SCAN_LIST>SCAN>PROCESSED"`
	Launch_datetime     string   `xml:"RESPONSE>SCAN_LIST>SCAN>LAUNCH_DATETIME"`
	Target              string   `xml:"RESPONSE>SCAN_LIST>SCAN>TARGET"`
	Type                string   `xml:"RESPONSE>SCAN_LIST>SCAN>TYPE"`
	Title               string   `xml:"RESPONSE>SCAN_LIST>SCAN>TITLE"`
	User_login          string   `xml:"RESPONSE>SCAN_LIST>SCAN>USER_LOGIN"`
	Ref                 string   `xml:"RESPONSE>SCAN_LIST>SCAN>REF"`
	Duration            string   `xml:"RESPONSE>SCAN_LIST>SCAN>DURATION"`
	State               string   `xml:"RESPONSE>SCAN_LIST>SCAN>STATUS>STATE"`
}

type PollScanOptions struct {
	Action  string `url:"action"` // list
	ScanRef string `url:"scan_ref"`
}

func (client *Client) PollScanResults(options *PollScanOptions) (*PollScanResponse, error) {
	options.Action = "list"

	urlString, err := addURLParameters(client.BaseURL.String(), options)

	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}

	var resp PollScanResponse

	_, requestError := client.MakeRequest(req, &resp)

	if requestError != nil {
		return nil, requestError
	}

	return &resp, nil
}

type CompletedScanResponse struct {
	XMLName                                                     xml.Name           `xml:"SIMPLE_RETURN"` // the actual outp was supposed to be: SCAN_LIST_OUTPUT?
	City                                                        string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>CITY"`
	Key                                                         []key              `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>KEY"`
	Name                                                        string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>NAME"`
	ZIPCode                                                     string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>ZIP_CODE"`
	Address                                                     string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>ADDRESS"`
	NameUserInfoHeaderComplianceScanResponse                    string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>USER_INFO>NAME"`
	Username                                                    string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>USER_INFO>USERNAME"`
	Role                                                        string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>USER_INFO>ROLE"`
	Hosts                                                       string             `xml:"RESPONSE>COMPLIANCE_SCAN>APPENDIX>TARGET_DISTRIBUTION>SCANNER>HOSTS"`
	Type                                                        string             `xml:"RESPONSE>COMPLIANCE_SCAN>APPENDIX>AUTHENTICATION>AUTH>TYPE"`
	Datetime                                                    string             `xml:"RESPONSE>DATETIME"`
	IP                                                          string             `xml:"RESPONSE>COMPLIANCE_SCAN>APPENDIX>AUTHENTICATION>AUTH>SUCCESS>IP"`
	State                                                       string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>STATE"`
	Country                                                     string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>COMPANY_INFO>COUNTRY"`
	HostsScanned                                                string             `xml:"RESPONSE>COMPLIANCE_SCAN>APPENDIX>TARGET_HOSTS>HOSTS_SCANNED"`
	NameScannerTargetDistributionAppendixComplianceScanResponse string             `xml:"RESPONSE>COMPLIANCE_SCAN>APPENDIX>TARGET_DISTRIBUTION>SCANNER>NAME"`
	GenerationDateTime                                          string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>GENERATION_DATETIME"`
	NameHeaderComplianceScanResponse                            string             `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>NAME"`
	OptionProfileTitle                                          optionProfileTitle `xml:"RESPONSE>COMPLIANCE_SCAN>HEADER>OPTION_PROFILE>OPTION_PROFILE_TITLE"`
}

type key struct {
	XMLName xml.Name `xml:"KEY"`
	Value   string   `xml:"value,attr"`
	Text    string   `xml:",chardata"`
}
type optionProfileTitle struct {
	XMLName                xml.Name `xml:"OPTION_PROFILE_TITLE"`
	Option_profile_default string   `xml:"option_profile_default,attr"`
	Text                   string   `xml:",chardata"`
}

type CompletedScanOptions struct {
	Action  string `url:"action"` // fetch
	ScanRef string `url:"scan_ref"`
}

func (client *Client) GetScanResults(options *CompletedScanOptions) (*CompletedScanResponse, error) {
	options.Action = "fetch"

	urlString, err := addURLParameters(client.BaseURL.String(), options)

	if err != nil {
		return nil, err
	}

	// i don't think there's any header data here?
	req, err := client.NewRequest(http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}

	var resp CompletedScanResponse

	_, requestError := client.MakeRequest(req, &resp)

	if requestError != nil {
		return nil, requestError
	}

	return &resp, nil
}
