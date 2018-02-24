package main

const (
	portCreateWhenCreatingInstance = `
  {
    "_context_roles": [
      "admin"
    ],
    "_context_request_id": "req-fdb23f2e-9c0e-46b1-802f-3194c1fad251",
    "event_type": "port.create.end",
    "timestamp": "2016-10-03 18:40:34.596836",
    "_context_tenant_id": "0b65cf220eab4a3cbd68681d188d7dc7",
    "_unique_id": "bca88f14c46e40559e981ac0b4ffebf5",
    "_context_tenant_name": "services",
    "_context_user": "31055c32b50442e5a4eb4c0f0cb3430b",
    "_context_user_id": "31055c32b50442e5a4eb4c0f0cb3430b",
    "payload": {
      "port": {
        "status": "DOWN",
        "binding:host_id": "oscomp-ch2-a06",
        "name": "",
        "allowed_address_pairs": [

        ],
        "admin_state_up": true,
        "network_id": "af33487a-4e96-4499-bfcd-4f741617a763",
        "tenant_id": "ada3b9b0dbac429f9361e803b54f5f32",
        "binding:vif_details": {
          "port_filter": true,
          "ovs_hybrid_plug": true
        },
        "binding:vnic_type": "normal",
        "binding:vif_type": "ovs",
        "device_owner": "compute:None",
        "mac_address": "fa:16:3e:4a:ac:75",
        "binding:profile": {
        },
        "fixed_ips": [
          {
            "subnet_id": "4a23cb36-b861-4daa-a8ef-c61360663669",
            "ip_address": "162.150.0.117"
          },
          {
            "subnet_id": "244c99a6-8011-4177-855b-dd493c5175c5",
            "ip_address": "2001:558:fe21:403:f816:3eff:fe4a:ac75"
          }
        ],
        "id": "a6c671d7-b4d5-4ebb-afaf-0c822bcc8948",
        "security_groups": [
          "0783a151-768c-49d3-a31d-178f70fabd51",
          "46d46540-98ac-4c93-ae62-68dddab2282e"
        ],
        "device_id": "128bc33a-22ae-48b4-8283-093b6ec749d0"
      }
    },
    "_context_project_name": "services",
    "_context_read_deleted": "no",
    "_context_tenant": "0b65cf220eab4a3cbd68681d188d7dc7",
    "priority": "INFO",
    "_context_is_admin": true,
    "_context_project_id": "0b65cf220eab4a3cbd68681d188d7dc7",
    "_context_timestamp": "2016-10-03 18:40:34.477012",
    "_context_user_name": "neutron",
    "publisher_id": "network.osctrl-ch2-a03",
    "message_id": "71047538-531f-4aca-be09-a31bec441d16"
  }

  `

	securityGroupRuleCreateWithCustomProtocall = `
  {  
     "_context_roles":[  
        "Member"
     ],
     "_context_request_id":"req-a17c784c-fec9-4077-8908-44b6f56b6196",
     "event_type":"security_group_rule.create.end",
     "timestamp":"2016-10-03 17:50:59.982008",
     "_context_tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_unique_id":"a7452605170c4979b2c6b76911d22026",
     "_context_tenant_name":"VOIP",
     "_context_user":"bca89c1b248e4aef9c69ece9e744cc54",
     "_context_user_id":"bca89c1b248e4aef9c69ece9e744cc54",
     "payload":{  
        "security_group_rule":{  
           "remote_group_id":null,
           "direction":"ingress",
           "protocol":10,
           "remote_ip_prefix":"10.0.0.0/8",
           "port_range_max":null,
           "dscp":null,
           "security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e",
           "tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
           "port_range_min":null,
           "ethertype":"IPv4",
           "id":"3eff38bb-eb03-450b-aed4-019d612baeec"
        }
     },
     "_context_project_name":"VOIP",
     "_context_read_deleted":"no",
     "_context_tenant":"ada3b9b0dbac429f9361e803b54f5f32",
     "priority":"INFO",
     "_context_is_admin":false,
     "_context_project_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_context_timestamp":"2016-10-03 17:50:59.925462",
     "_context_user_name":"admin",
     "publisher_id":"network.osctrl-ch2-a03",
     "message_id":"6c93e24f-0892-494b-8e68-46252ceb9611"
  }
  `

	securityGroupRuleCreateWithIcmpAndCider = `
  {  
     "_context_roles":[  
        "Member"
     ],
     "_context_request_id":"req-c584fd21-9e58-4624-b316-b53487eed98e",
     "event_type":"security_group_rule.create.end",
     "timestamp":"2016-10-03 18:05:35.836029",
     "_context_tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_unique_id":"cd280fd4f1474266bd0ad6e3ee5933a6",
     "_context_tenant_name":"VOIP",
     "_context_user":"bca89c1b248e4aef9c69ece9e744cc54",
     "_context_user_id":"bca89c1b248e4aef9c69ece9e744cc54",
     "payload":{  
        "security_group_rule":{  
           "remote_group_id":null,
           "direction":"ingress",
           "protocol":"icmp",
           "remote_ip_prefix":"192.168.1.0/24",
           "port_range_max":null,
           "dscp":null,
           "security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e",
           "tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
           "port_range_min":null,
           "ethertype":"IPv4",
           "id":"66d7ac79-3551-4436-83c7-103b50760cfb"
        }
     },
     "_context_project_name":"VOIP",
     "_context_read_deleted":"no",
     "_context_tenant":"ada3b9b0dbac429f9361e803b54f5f32",
     "priority":"INFO",
     "_context_is_admin":false,
     "_context_project_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_context_timestamp":"2016-10-03 18:05:35.769947",
     "_context_user_name":"admin",
     "publisher_id":"network.osctrl-ch2-a03",
     "message_id":"f67b70d5-a782-4c5e-a274-a7ff197b73ec"
  }
  `
	securityGroupRuleCreateWithports = `
  {  
     "_context_roles":[  
        "Member"
     ],
     "_context_request_id":"req-1f17d667-c33f-4fa4-a026-8e2872dbf1d8",
     "event_type":"security_group_rule.create.end",
     "timestamp":"2016-10-03 17:32:25.723344",
     "_context_tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_unique_id":"2fad8ecdd86e4748850d91bb0c83d625",
     "_context_tenant_name":"VOIP",
     "_context_user":"bca89c1b248e4aef9c69ece9e744cc54",
     "_context_user_id":"bca89c1b248e4aef9c69ece9e744cc54",
     "payload":{  
        "security_group_rule":{  
           "remote_group_id":null,
           "direction":"ingress",
           "protocol":"tcp",
           "remote_ip_prefix":"10.0.0.0/8",
           "port_range_max":443,
           "dscp":null,
           "security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e",
           "tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
           "port_range_min":443,
           "ethertype":"IPv4",
           "id":"2b84d898-67b4-4370-9808-40a3fdb55a64"
        }
     },
     "_context_project_name":"VOIP",
     "_context_read_deleted":"no",
     "_context_tenant":"ada3b9b0dbac429f9361e803b54f5f32",
     "priority":"INFO",
     "_context_is_admin":false,
     "_context_project_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_context_timestamp":"2016-10-03 17:32:25.665588",
     "_context_user_name":"admin",
     "publisher_id":"network.osctrl-ch2-a03",
     "message_id":"4df01871-8bdb-4b85-bb34-cbff59ee6034"
  }
  `
	securityGroupRuleCreateWithSecurityGroupAsRemoteIPPrefix = `
  {  
     "_context_roles":[  
        "Member"
     ],
     "_context_request_id":"req-9e0360c7-786f-4a5b-84b6-7d2ccd23cbdd",
     "event_type":"security_group_rule.create.end",
     "timestamp":"2016-10-03 17:36:58.780554",
     "_context_tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_unique_id":"b38fe8caed514eb2ba910e1ae74c6321",
     "_context_tenant_name":"VOIP",
     "_context_user":"bca89c1b248e4aef9c69ece9e744cc54",
     "_context_user_id":"bca89c1b248e4aef9c69ece9e744cc54",
     "payload":{  
        "security_group_rule":{  
           "remote_group_id":"0783a151-768c-49d3-a31d-178f70fabd51",
           "direction":"ingress",
           "protocol":"tcp",
           "remote_ip_prefix":null,
           "port_range_max":25,
           "dscp":null,
           "security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e",
           "tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
           "port_range_min":20,
           "ethertype":"IPv6",
           "id":"7b14b6cd-f966-4b61-aaad-c03d8eacc830"
        }
     },
     "_context_project_name":"VOIP",
     "_context_read_deleted":"no",
     "_context_tenant":"ada3b9b0dbac429f9361e803b54f5f32",
     "priority":"INFO",
     "_context_is_admin":false,
     "_context_project_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_context_timestamp":"2016-10-03 17:36:58.712962",
     "_context_user_name":"admin",
     "publisher_id":"network.osctrl-ch2-a03",
     "message_id":"e2d7c089-8194-4523-8f84-ae22db497f60"
  }
  `
	securityGroupRuleCreateWithSSHOpenToTheInternet = `
  {  
     "_context_roles":[  
        "Member"
     ],
     "_context_request_id":"req-94df69c6-1c3f-48bd-b2f6-f47abdef5d9b",
     "event_type":"security_group_rule.create.end",
     "timestamp":"2016-10-03 18:09:11.938476",
     "_context_tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_unique_id":"09412fff881543679f30412ef2342954",
     "_context_tenant_name":"VOIP",
     "_context_user":"bca89c1b248e4aef9c69ece9e744cc54",
     "_context_user_id":"bca89c1b248e4aef9c69ece9e744cc54",
     "payload":{  
        "security_group_rule":{  
           "remote_group_id":null,
           "direction":"ingress",
           "protocol":"tcp",
           "remote_ip_prefix":"0.0.0.0/0",
           "port_range_max":22,
           "dscp":null,
           "security_group_id":"46d46540-98ac-4c93-ae62-68dddab2282e",
           "tenant_id":"ada3b9b0dbac429f9361e803b54f5f32",
           "port_range_min":22,
           "ethertype":"IPv4",
           "id":"bf288dfc-f9cb-446b-bacc-a8de016c9b11"
        }
     },
     "_context_project_name":"VOIP",
     "_context_read_deleted":"no",
     "_context_tenant":"ada3b9b0dbac429f9361e803b54f5f32",
     "priority":"INFO",
     "_context_is_admin":false,
     "_context_project_id":"ada3b9b0dbac429f9361e803b54f5f32",
     "_context_timestamp":"2016-10-03 18:09:11.876789",
     "_context_user_name":"admin",
     "publisher_id":"network.osctrl-ch2-a03",
     "message_id":"afb043b6-fa56-470b-b17e-984fb4cb6505"
  }
  `
	securityGroupRuleDeleteWithIcmpAndCider = `
  {
    "_context_roles": [
      "Member"
    ],
    "_context_request_id": "req-836eb80f-c6eb-459b-87b6-a093ebac3051",
    "event_type": "security_group_rule.delete.end",
    "timestamp": "2016-10-03 18:14:33.007074",
    "_context_tenant_id": "ada3b9b0dbac429f9361e803b54f5f32",
    "_unique_id": "04beeb34769b43bca09ec837d86ed18b",
    "_context_tenant_name": "VOIP",
    "_context_user": "bca89c1b248e4aef9c69ece9e744cc54",
    "_context_user_id": "bca89c1b248e4aef9c69ece9e744cc54",
    "payload": {
      "security_group_rule_id": "7b14b6cd-f966-4b61-aaad-c03d8eacc830"
    },
    "_context_project_name": "VOIP",
    "_context_read_deleted": "no",
    "_context_tenant": "ada3b9b0dbac429f9361e803b54f5f32",
    "priority": "INFO",
    "_context_is_admin": false,
    "_context_project_id": "ada3b9b0dbac429f9361e803b54f5f32",
    "_context_timestamp": "2016-10-03 18:14:32.962116",
    "_context_user_name": "admin",
    "publisher_id": "network.osctrl-ch2-a03",
    "message_id": "9bc5106c-a08b-4cda-9311-20bc16bc3008"
  }
  `
)
