package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestProcessWaitingEvent(t *testing.T) {
	var delivery amqp.Delivery
	openstackActions := connectFakeOpenstack()

	delivery.Body = []byte(securityGroupRuleCreateWithIcmpAndCider)
	event, err := processWaitingEvent(delivery, openstackActions)
	if err != nil {
		t.Fatal(err)
	}

	_ = event
}

func TestLogEvents(t *testing.T) {
	hostName, _ := os.Hostname()
	IPList := []string{"10.0.0.1", "10.0.0.3"}
	logLines := []string{fmt.Sprintf(`{"security_group_rule":{"remote_group_id":null,"direction":"ingress","protocol":"icmp","remote_ip_prefix":"192.168.1.0/24","port_range_max":null,"rule_direction":"","security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e","tenant_id":"ada3b9b0dbac429f9361e803b54f5f32","port_range_min":null,"ethertype":"IPv4","id":"66d7ac79-3551-4436-83c7-103b50760cfb"},"affected_ip_address":"10.0.0.1","change_type":"sg_rule_add","source_type":"osel","source_message_bus":"%s"}`, hostName), fmt.Sprintf(`{"security_group_rule":{"remote_group_id":null,"direction":"ingress","protocol":"icmp","remote_ip_prefix":"192.168.1.0/24","port_range_max":null,"rule_direction":"","security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e","tenant_id":"ada3b9b0dbac429f9361e803b54f5f32","port_range_min":null,"ethertype":"IPv4","id":"66d7ac79-3551-4436-83c7-103b50760cfb"},"affected_ip_address":"10.0.0.3","change_type":"sg_rule_add","source_type":"osel","source_message_bus":"%s"}`, hostName)}
	logger := connectFakeSyslog()
	qualys := connectFakeQualys()
	IPs := make(map[string][]string)

	IPs["46d46540-98ac-4c93-ae62-68dddab2282e"] = IPList
	fakeEvent := Event{
		RawData:   []byte(securityGroupRuleCreateWithIcmpAndCider),
		LogLines:  logLines,
		Processor: EventSecurityGroupRuleChange{ChangeType: "sg_rule_add"},
		IPs:       IPs,
	}
	events := []Event{fakeEvent}

	logEvents(events, logger, qualys)
	savedLogs := logger.GetLogs()
	assert.Equal(t, 2, len(savedLogs))

	logLine1 := fmt.Sprintf(`{"affected_ip_address":"10.0.0.1","change_type":"sg_rule_add","qualys_scan_id":"","qualys_scan_error":"Not scanned by Qualys","security_group_rule":{"remote_group_id":null,"direction":"ingress","protocol":"icmp","remote_ip_prefix":"192.168.1.0/24","port_range_max":null,"rule_direction":"","security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e","tenant_id":"ada3b9b0dbac429f9361e803b54f5f32","port_range_min":null,"ethertype":"IPv4","id":"66d7ac79-3551-4436-83c7-103b50760cfb"},"source_type":"osel1.1","source_message_bus":"%s"}`, hostName)
	assert.Equal(t, logLine1, savedLogs[0])
}
