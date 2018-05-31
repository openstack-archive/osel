# OpenStack Event Listener

## What is STORY.md?

Retweeted Safia Abdalla (@captainsafia):

From now on, you can expect to see a "STORY.md" in each of my @github repos
that describes the technical story/process behind the project.

https://twitter.com/captainsafia/status/839587421247389696

## Introduction

I wanted to write a little about a project that I enjoyed working on, called
the OpenStack Event Listener, or "OSEL" for short.  This project bridged the
OpenStack control plane on one hand, and an external scanning facility provided
by Qualys.  It had a number of interesting challenges.  I was never able to
really concentrate on it - this project took about 20% of my time for a period
of about 3 months.

I am writing this partially as catharsis, to allow my brain to mark this part
of my mental inventory as ripe for reclamation.  I am also writing on the off
chance that someone might find this useful in the future.

## The Setting

Let me paint a picture of the environment in which this development occurred.  

The Comcast OpenStack environment was transitioning from the OpenStack Icehouse
release (very old) to the Newton release (much more current).  This development
occurred within the context of the Icehouse environment.  

Comcast's security team uses S3 RiskFabric to manage auditing and tracking
security vulnerabilities across the board.  They also engage the services of
Qualys to perform network scanning (in a manner very similar to Nessus) once a
day against all the CIDR blocks that comprise Comcast's Internet-routable IP
addresses.  Qualys scanning could also be triggered on-demand.

## Technical Requirements

First, let me describe the technical requirements for OSEL:

* OSEL would connect to the OpenStack RabbitMQ message bus and register as a
  listener for "notification" events.  This would allow OSEL to inspect all
  events, including security group changes.
* When a security group change occurred, OSEL would ensure that it had the
  details of the change (ports permitted or blocked) as well as a list of all
  affected IP addresses.
* OSEL would initiate a Qualys scan using the Qualys API.  This would return a
  scan ID.
* OSEL would log the change as well as the Qualys scan ID to the Security
  instance of Splunk to create an audit trail.
* Qualys scan results would be imported into S3 RiskFabric for security audit
  management.

## Implementation Approach

My group does most of it's development in Go, and this was no exception.  

This is what the data I was getting back from the AMQP message looked like.
All identifiers have been scrambled.

```json
{  
   "_context_roles":[  
      "Member"
   ],
   "_context_request_id":"req-f96ea9a5-435e-4177-8e51-bfe60d0fae2a",
   "event_type":"security_group_rule.create.end",
   "timestamp":"2016-10-03 18:10:59.112712",
   "_context_tenant_id":"ada3b9b06482909f9361e803b54f5f32",
   "_unique_id":"eafc9362327442b49d8c03b0e88d0216",
   "_context_tenant_name":"BLURP",
   "_context_user":"bca89c1b248e4a78282899ece9e744cc54",
   "_context_user_id":"bca89c1b248e4a78282899ece9e744cc54",
   "payload":{  
      "security_group_rule_id":"bf8318fc-f9cb-446b-ffae-a8de016c562"
   },
   "_context_project_name":"BLURP",
   "_context_read_deleted":"no",
   "_context_tenant":"ada3b9b06482909f9361e803b54f5f32",
   "priority":"INFO",
   "_context_is_admin":false,
   "_context_project_id":"ada3b9b06482909f9361e803b54f5f32",
   "_context_timestamp":"2016-10-03 18:10:59.079179",
   "_context_user_name":"admin",
   "publisher_id":"network.osctrl1",
   "message_id":"e75fb2ee-85bf-44ba-a083-2445eca2ae10"
}
```

## Testing Pattern

I leaned heavily on dependency injection to make this code as testable as
possible.  For example, I needed an object that would contain the persistent
`syslog.Writer`.  I created a `SyslogActioner` interface to represent all
interactions with syslog.  When the code is operating normally, interactions
with syslog occur through methods of the `SyslogActions` struct.  But in unit
testing mode the `SyslogTestActions` struct is used instead, and all that does
is save copies of all messages that would have been sent so they can be
compared against the intended messages.  This facilitates good testing.

## Fate of the Project

The OSEL project was implemented and installed into production.  There was a
problem with it.

There was no exponential backoff for the AMQP connection to the OpenStack
control plane's RabbitMQ.  When that RabbitMQ had issues - which was
surprisingly often - OSEL would hanner away, trying to connect to it.  That
would not be too much of an issue; despite what was effectively an infinite
loop, CPU usage was not extreme.  The real problem was that connection failures
were logged - and logs could become several gigabytes in a matter of hours.
This was mitigated by the OpenStack operations team rotating the logs hourly,
and alerting if an hour's worth of logs exceeded a set size.  It was my
intention to use one of the many [exponential backoff
modules](https://github.com/cenkalti/backoff) available out there to make this
more graceful.

## Remaining Work

If OSEL were ever to be un-shelved, here are a few of the things that I wish I
had time to implement.

- Neutron Port Events: The initial release of OSEL processed only security
  group rule additions, modifications, or deletions.  So that covered the base
  case for when a security group was already associated with a set of OpenStack
  Networking (neutron) ports.   But a scan should be similarly launched when a
  new port is created and associated to a security group.  This is what happens
  when a new host is created.
- Modern OpenStack: In order to make this work with a more modern OpenStack, it
  would probably best to integrate with events generated through Aodh.  Aodh
  seems to be built for this kind of reporting. 
- Implement exponential backoff for AMQP connections as mentioned earlier.
