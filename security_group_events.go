package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// EventSecurityGroupRuleChange is the event processor class for all changes to
// security groups.  This includes additions and deletions.  This must conform
// to the EventProcessor interface (see events.go).
type EventSecurityGroupRuleChange struct {
	ChangeType string
}

// FillExtraData takes a security group change and enriches it with additional
// information about the affected IP addresses using the
// OpenStackActionInterface getPortList function.
func (s EventSecurityGroupRuleChange) FillExtraData(e *Event, openstack OpenStackActioner) error {
	// PopulateIps: This function returns a map of security group to array of IP addresses for all ports in the specified tenantID.

	err := openstack.Connect(e.EventData.TenantID, e.EventData.UserName)
	if err != nil {
		return err
	}

	// Make port list request to neutron
	resultMap := openstack.GetPortList()
	resultIPAddresses := make(map[string][]string)
	for _, ipMap := range resultMap {
		resultIPAddresses[ipMap.securityGroup] = append(resultIPAddresses[ipMap.securityGroup], ipMap.ipAddress)
	}
	e.IPs = resultIPAddresses
	return nil
}

// FormatLogs takes the accumulated event data and composes the JSON message to
// be logged.
func (s EventSecurityGroupRuleChange) FormatLogs(e *Event, scannedIPAddresses []string) ([]string, error) {
	var es osSecurityGroupRuleChange
	var logLines []string
	if e == nil {
		return logLines, fmt.Errorf("Event must not be nil")
	}

	if err := json.Unmarshal(e.RawData, &es); err != nil {
		return logLines, err
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	es.Payload.ChangeType = s.ChangeType
	es.Payload.SourceType = OselVersion
	es.Payload.SourceMessageBus = hostName
	es.Payload.QualysScanID = e.QualysScanID
	es.Payload.QualysScanError = e.QualysScanError

	affectedIPArray := e.IPs[es.Payload.SecurityGroupRule.SecurityGroupID]
	qualysScanJoin := fmt.Sprintf("|%s|", strings.Join(scannedIPAddresses, "|"))
	for _, affectedIPAddr := range affectedIPArray {
		es.Payload.QualysScanID = ""
		es.Payload.QualysScanError = ""
		if strings.Index(qualysScanJoin, fmt.Sprintf("|%s|", affectedIPAddr)) > -1 {
			es.Payload.QualysScanID = e.QualysScanID
			es.Payload.QualysScanError = e.QualysScanError
		} else {
			es.Payload.QualysScanID = ""
			es.Payload.QualysScanError = "Not scanned by Qualys"
		}
		es.Payload.AffectedIPAddr = affectedIPAddr
		jsonLine, err := json.Marshal(es.Payload)
		if err != nil {
			return nil, err
		}
		logLines = append(logLines, string(jsonLine))
	}
	return logLines, nil
}
