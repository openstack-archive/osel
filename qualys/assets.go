package qualys

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const (
	assetsBasePath = "asset"
	groupsBasePath = "group"
)

// AssetsService is an interface for interfacing with the Assets
// endpoints of the Qualys API
type AssetsService interface {
	ListAssetGroups(*ListAssetGroupOptions) ([]AssetGroup, *Response, error)
	GetAssetGroupByID(groupID string) (*AssetGroup, *Response, error)
	AddIPsToGroup(*AddIPsToGroupOptions) (*Response, error)
}

// AssetsServiceOp handles communication with the asset related methods of the
// Qualys API.
type AssetsServiceOp struct {
	client *Client
}

var _ AssetsService = &AssetsServiceOp{}

// AssetGroup represents a Qualys HostGroup
type AssetGroup struct {
	ID          string        `xml:"ID"`
	Title       string        `xml:"TITLE"`
	OwnerUserID string        `xml:"OWNER_USER_ID"`
	OwnerUnitID string        `xml:"OWNER_UNIT_ID"`
	IPs         AssetGroupIPs `xml:"IP_SET"`
}

// AssetGroupIPs represents one or more IP addresses assigned to the AssetGroup
type AssetGroupIPs struct {
	IPs      []string `xml:"IP"`
	IPRanges []string `xml:"IP_RANGE"`
}

type ipRange struct {
	Min net.IP
	Max net.IP
}

func newIPRange(rangeString string) *ipRange {
	var r = strings.Split(rangeString, "-")
	return &ipRange{Min: net.ParseIP(r[0]), Max: net.ParseIP(r[1])}
}

func (ip *ipRange) Contains(ipString string) bool {
	var myIP = net.ParseIP(ipString)
	if bytes.Compare(myIP, ip.Min) >= 0 && bytes.Compare(myIP, ip.Max) <= 0 {
		return true
	}
	return false
}

// ContainsIP returns true when the AssetGroupIPs matches the provided IP
func (agp *AssetGroupIPs) ContainsIP(ip string) bool {
	if containsString(agp.IPs, ip) {
		return true
	}
	if agp.IPRanges != nil && len(agp.IPRanges) > 0 {
		for _, ipRange := range agp.IPRanges {
			if newIPRange(ipRange).Contains(ip) {
				return true
			}
		}
	}
	return false
}

// ContainsIP returns true when the AssetGroup has any assets matching the provided IP
func (ag *AssetGroup) ContainsIP(ip string) bool {
	return ag.IPs.ContainsIP(ip)
}

type assetGroupsRoot struct {
	AssetGroups []AssetGroup `xml:"RESPONSE>ASSET_GROUP_LIST>ASSET_GROUP"`
}

// AssetGroupUpdateRequest represents a request to update a group
type AssetGroupUpdateRequest struct {
}

// ListAssetGroupOptions represents the AssetGroup retrieval options
type ListAssetGroupOptions struct {
	Ids    []string `url:"ids,comma,omitempty"`
	Action string   `url:"action,omitempty"`
}

// AddIPsToGroupOptions represents the update request for an AssetGroup
type AddIPsToGroupOptions struct {
	GroupID string   `url:"id,omitempty"`
	IPs     []string `url:"add_ips,comma,omitempty"`
}

// ListAssetGroups retrieves a list of AssetGroups
func (s *AssetsServiceOp) ListAssetGroups(opt *ListAssetGroupOptions) ([]AssetGroup, *Response, error) {
	return s.listAssetGroups(opt)
}

// GetAssetGroupByID retrieves an AssetGroup by id.
func (s *AssetsServiceOp) GetAssetGroupByID(groupID string) (*AssetGroup, *Response, error) {
	return s.getAssetGroup(groupID)
}

// AddIPsToGroup adds the IPs in AddIPsToGroupOptions to the AssetGroup
func (s *AssetsServiceOp) AddIPsToGroup(opt *AddIPsToGroupOptions) (*Response, error) {
	return s.addIPsToGroup(opt)
}

func (s *AssetsServiceOp) getAssetGroup(groupID string) (*AssetGroup, *Response, error) {
	opts := ListAssetGroupOptions{Ids: []string{groupID}}
	groups, response, err := s.listAssetGroups(&opts)
	if err != nil {
		return nil, response, err
	}
	if len(groups) == 0 {
		return nil, response, nil
	}
	return &groups[0], response, nil
}

func (s *AssetsServiceOp) addIPsToGroup(opt *AddIPsToGroupOptions) (*Response, error) {
	path := fmt.Sprintf("%s/%s/?action=edit", assetsBasePath, groupsBasePath)

	req, err := s.client.NewRequest("POST", path, opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.MakeRequest(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, err
}

// Helper method for listing asset groups
func (s *AssetsServiceOp) listAssetGroups(listOpt *ListAssetGroupOptions) ([]AssetGroup, *Response, error) {
	path := fmt.Sprintf("%s/%s/", assetsBasePath, groupsBasePath)
	if listOpt == nil {
		listOpt = &ListAssetGroupOptions{}
	}
	listOpt.Action = "list"
	path, err := addURLParameters(path, listOpt)

	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(assetGroupsRoot)
	resp, err := s.client.MakeRequest(req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.AssetGroups, resp, err
}
