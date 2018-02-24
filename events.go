package main

import (
	"encoding/json"
	"log"
	"strings"
)

// EventProcessor is an Interface for event-specific classes that will process
// events based on their specific fiends.
type EventProcessor interface {
	FormatLogs(*Event, []string) ([]string, error)
	FillExtraData(*Event, OpenStackActioner) error
}

// Event is a class representing an event accepted from the AMQP, and the
// additional attributes that have been parsed from it.
type Event struct {
	EventData          *openStackEvent
	RawData            []byte
	IPs                map[string][]string
	SecurityGroupRules []*osSecurityGroupRule
	LogLines           []string
	Processor          EventProcessor
	QualysScanID       string
	QualysScanError    string
}

// ParseEvent takes the []byte that has been received from the AMQP message,
// demarshals the JSON, and then returns the event data as well as an event
// processor specific to that type of event.
func ParseEvent(message []byte) (Event, error) {
	var osEvent openStackEvent
	if err := json.Unmarshal(message, &osEvent); err != nil {
		return Event{}, err
	}

	e := Event{
		EventData: &osEvent,
		RawData:   message,
	}

	if Debug {
		log.Printf("Event detected: %s\n", osEvent.EventType)
	}

	switch {
	case strings.Contains(e.EventData.EventType, "security_group_rule.create.end"):
		e.Processor = EventSecurityGroupRuleChange{ChangeType: "sg_rule_add"}
	case strings.Contains(e.EventData.EventType, "security_group_rule.delete.end"):
		e.Processor = EventSecurityGroupRuleChange{ChangeType: "sg_rule_del"}
		// case strings.Contains(e.EventData.EventType, "port.create.end"):
		// 	e.Processor = EventPortChange{ChangeType: "port_create"}
	}

	return e, nil
}
