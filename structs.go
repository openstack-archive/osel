package main

type openStackEvent struct {
	EventType   string `json:"event_type"`
	Timestamp   string `json:"timestamp"`
	TenantID    string `json:"_context_tenant_id"`
	TenantName  string `json:"_context_tenant_name"`
	User        string `json:"_context_user"`
	UserName    string `json:"_context_user_name"`
	UserID      string `json:"_context_user_id"`
	IsAdmin     bool   `json:"_context_is_admin"`
	PublisherID string `json:"publisher_id"`
	MessageID   string `json:"message_id"`
}

type osSecurityGroupRule struct {
	RemoteGroupID  interface{} `json:"remote_group_id"`
	Direction      string      `json:"direction"`
	Protocol       interface{} `json:"protocol"`
	RemoteIPPrefix string      `json:"remote_ip_prefix"`
	PortRangeMax   interface{} `json:"port_range_max"`
	//	Dscp            interface{} `json:"dscp"`
	Rule            string      `json:"rule_direction"`
	SecurityGroupID string      `json:"security_group_id"`
	TenantID        string      `json:"tenant_id"`
	PortRangeMin    interface{} `json:"port_range_min"`
	Ethertype       string      `json:"ethertype"`
	ID              string      `json:"id"`
}

type osSecurityGroupRuleChange struct {
	Payload struct {
		AffectedIPAddr    interface{}         `json:"affected_ip_address"`
		ChangeType        string              `json:"change_type"`
		QualysScanID      string              `json:"qualys_scan_id"`
		QualysScanError   string              `json:"qualys_scan_error"`
		SecurityGroupRule osSecurityGroupRule `json:"security_group_rule"`
		SourceType        string              `json:"source_type"`
		SourceMessageBus  string              `json:"source_message_bus"`
	} `json:"payload"`
}

type osSecurityGroupRuleDelete struct {
	Payload struct {
		SecurityGroupRuleID string `json:"security_group_rule_id"`
	} `json:"payload"`
}

type osPortCreate struct {
	Payload struct {
		Port osPort `json:"port"`
	} `json:"payload"`
}

type osPort struct {
	Status              string        `json:"status"`
	BindingHostID       string        `json:"binding:host_id"`
	Name                string        `json:"name"`
	AllowedAddressPairs []interface{} `json:"allowed_address_pairs"`
	AdminStateUp        bool          `json:"admin_state_up"`
	NetworkID           string        `json:"network_id"`
	TenantID            string        `json:"tenant_id"`
	BindingVifDetails   struct {
		PortFilter    bool `json:"port_filter"`
		OvsHybridPlug bool `json:"ovs_hybrid_plug"`
	} `json:"binding:vif_details"`
	BindingVnicType string `json:"binding:vnic_type"`
	BindingVifType  string `json:"binding:vif_type"`
	DeviceOwner     string `json:"device_owner"`
	MacAddress      string `json:"mac_address"`
	BindingProfile  struct {
	} `json:"binding:profile"`
	FixedIps []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
	ID             string   `json:"id"`
	SecurityGroups []string `json:"security_groups"`
	DeviceID       string   `json:"device_id"`
}
