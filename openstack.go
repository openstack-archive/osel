package main

/*

openstack - This file includes all of the logic necessary to interact with
OpenStack.  This is extrapolated out so that an OpenStackActioner
interface can be passed to functions.  Doing this allows testing by mock
classes to be created that can be passed to functions.

Since this is a wrapper around the gophercloud libraries, this does not need
testing.

*/

import (
	"fmt"
	"log"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
)

// OpenStackActioner is an interface for an OpenStackActions class.
// Having this as an interface allows us to pass in a dummy class for testing
// that just returns mocked data.
type OpenStackActioner interface {
	GetPortList() []OpenStackIPMap
	Connect(string, string) error
}

// OpenStackActions is a class that handles all interactions directly with
// OpenStack.  See the comment on OpenStackActioner for rationale.
type OpenStackActions struct {
	gopherCloudClient *gophercloud.ProviderClient
	neutronClient     *gophercloud.ServiceClient
	Options           OpenStackOptions
}

// OpenStackOptions is a class to convey all of the configurable options for the
// OpenStackActions class.
type OpenStackOptions struct {
	KeystoneURI string
	Password    string
	RegionName  string
	TenantID    string
	UserName    string
}

// OpenStackIPMap is a struct that is used to capture the mapping of IP address
// to security group.  It is what is returned, in array form, from port list.
type OpenStackIPMap struct {
	ipAddress     string
	securityGroup string
}

// GetPortList is a method that uses GopherCloud to query OpenStack for a
// list of ports, with their associated security group. It returns an array of
// OpenStackIPMap.
func (s OpenStackActions) GetPortList() []OpenStackIPMap {
	// Make port list request to neutron
	var ips []OpenStackIPMap
	portListOpts := ports.ListOpts{
		TenantID: s.Options.TenantID,
	}
	if s.neutronClient == nil {
		log.Println("Error: neutronClient is nil")
	}
	pager := ports.List(s.neutronClient, portListOpts)

	// Define an anonymous function to be executed on each page's iteration
	pager.EachPage(func(page pagination.Page) (bool, error) {
		portList, err := ports.ExtractPorts(page)
		if err != nil {
			// ignore ?
		}

		for _, p := range portList {
			// "p" will be a ports.Port
			for _, fixedIP := range p.FixedIPs {
				for _, securityGroup := range p.SecurityGroups {
					ips = append(ips, OpenStackIPMap{
						ipAddress:     fixedIP.IPAddress,
						securityGroup: securityGroup,
					})
				}
			}
		}
		return true, err
	})
	return ips
}

// Connect is the method that establishes a connection to the OpenStack
// service.
func (s *OpenStackActions) Connect(tenantID string, username string) error {
	var err error
	keystoneOpts := gophercloud.AuthOptions{
		IdentityEndpoint: s.Options.KeystoneURI,
		TenantID:         tenantID,
		Username:         username,
		Password:         s.Options.Password,
		AllowReauth:      true,
	}

	log.Println(fmt.Sprintf("Connecting to keystone %q in region %q for tenant %q with user %q", s.Options.KeystoneURI,
		s.Options.RegionName, tenantID, username))
	s.gopherCloudClient, err = openstack.AuthenticatedClient(keystoneOpts)
	if err != nil {
		return fmt.Errorf("unable to connect to %s using user %s for tenant %s: %s",
			s.Options.KeystoneURI, s.Options.UserName, s.Options.TenantID, err)
	}
	log.Println("Connected to gophercloud ", s.Options.KeystoneURI)

	neutronOpts := gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: s.Options.RegionName,
	}
	s.neutronClient, err = openstack.NewNetworkV2(s.gopherCloudClient, neutronOpts)
	if err != nil {
		return fmt.Errorf("unable to connect to neutron using user %s in region %s: %s",
			s.Options.UserName, s.Options.RegionName, err)
	}
	return nil
}
