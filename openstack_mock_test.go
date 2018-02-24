package main

type OpenStackTestActions struct {
	regionName string
	tenantID   string
}

func (s OpenStackTestActions) Connect(tenantID string, username string) error {
	return nil
}

func (s OpenStackTestActions) GetPortList() []OpenStackIPMap {
	return []OpenStackIPMap{
		{
			ipAddress:     "10.0.0.1",
			securityGroup: "46d46540-98ac-4c93-ae62-68dddab2282e",
		},
		{
			ipAddress:     "10.0.0.2",
			securityGroup: "groupTwo",
		},
		{
			ipAddress:     "10.0.0.3",
			securityGroup: "46d46540-98ac-4c93-ae62-68dddab2282e",
		},
	}
}

func connectFakeOpenstack() *OpenStackTestActions {
	return new(OpenStackTestActions)
}
