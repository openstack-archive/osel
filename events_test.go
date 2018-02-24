package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEventWillReturnAnEventStruct(t *testing.T) {
	event, err := ParseEvent([]byte(securityGroupRuleCreateWithports))
	assert.Nil(t, err)
	assert.Equal(t, "main.Event", reflect.TypeOf(event).String(),
		"ParseEvent should return an Event struct")
	assert.Equal(t, "security_group_rule.create.end", event.EventData.EventType)
	assert.Equal(t, "bca89c1b248e4aef9c69ece9e744cc54", event.EventData.UserID)
	assert.Equal(t, "admin", event.EventData.UserName)
	assert.Equal(t, "ada3b9b0dbac429f9361e803b54f5f32", event.EventData.TenantID)
	assert.Equal(t, "VOIP", event.EventData.TenantName)
}

func TestParseEventWillCreateTheProperEventProcessor(t *testing.T) {
	e, err := ParseEvent([]byte(securityGroupRuleCreateWithports))
	assert.Nil(t, err)
	//assert.Equal(t, "main.EventSecurityGroupRuleChange", reflect.TypeOf(e.Processor).String(),
	//	"ParseEvent should return the proper implementation of EventProcessor")
	assert.Equal(t, EventSecurityGroupRuleChange{"sg_rule_add"}, e.Processor,
		"ParseEvent should return the proper implementation of EventProcessor")

	e, err = ParseEvent([]byte(securityGroupRuleDeleteWithIcmpAndCider))
	assert.Nil(t, err)
	assert.Equal(t, "main.EventSecurityGroupRuleChange", reflect.TypeOf(e.Processor).String(),
		"ParseEvent should return the proper implementation of EventProcessor")

	//	_, eventProcessor, err = ParseEvent([]byte(portCreateWhenCreatingInstance))
	//	assert.Nil(t, err)
	//	assert.Equal(t, "main.EventPortChange", reflect.TypeOf(eventProcessor).String(),
	//		"ParseEvent should return the proper implementation of EventProcessor")

}

// func TestPortCreateEvent(t *testing.T) {
//	fakeOpenStack := connectFakeOpenstack()
//	event, eventProcessor, err := ParseEvent([]byte(portCreateWhenCreatingInstance))
//	assert.Nil(t, err)
//	eventProcessor.FillExtraData(&event, fakeOpenStack)
//}

func TestEventSecurityGroupRuleCreateEvent(t *testing.T) {
	fakeOpenStack := connectFakeOpenstack()
	event, err := ParseEvent([]byte(securityGroupRuleCreateWithports))
	assert.Nil(t, err)
	event.Processor.FillExtraData(&event, fakeOpenStack)
}

func TestEventSecurityGroupRuleDeleteEvent(t *testing.T) {
	fakeOpenStack := connectFakeOpenstack()
	event, err := ParseEvent([]byte(securityGroupRuleDeleteWithIcmpAndCider))
	assert.Nil(t, err)
	event.Processor.FillExtraData(&event, fakeOpenStack)
}
