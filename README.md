OpenStack Event Listener
========================

What does this do?
------------------

The OpenStack Event Listener connects to the OpenStack message bus (RabbitMQ)
and listens for certain kinds of events.  When it detects those events, it will
gather additional data and forward the information to external systems for
processing.  It integrates with syslog and the Qualys API.

The initial use case that inspired this project was to detect when security
group changes occurred and to trigger an external port scan of the affected IP
addresses so that we could ensure that the change did not create a new
vulnerability by opening something up to the Internet.

For more background information on this project, see [the story of
osel](STORY.md).

Current State
-------------
Code maturity is considered experimental.

Installation
------------
Use `go get git.openstack.org/openstack/osel`.  Or alternatively,
download or clone the repository.

The lib was developed and tested on go 1.10. 

Configuration
-------------

Configuration resides in a YAML-format configuration file.  Before running the
os_event_listener process set the EL_CONFIG environment variable to the
absolute path to that file.

This is an example of the configuration format:

```yaml
debug: true
batch_interval: 2
rabbit_uri: "amqp://amqp_user:amqp_password@amqp_host:amqp_port//"
logfile: "/var/log/os_event_listener.log"
syslog_server: your.syslog.server.fqdn
syslog_port: "514"
syslog_protocol: "tcp"
retry_syslog: "false"
openstack:
  identity_endpoint: "https://keystone.url:5000/v2.0/"
  tenant_name: "tenant_to_authenticate_against"
  user: "username"
  password: "password"
  region: "region_name"
qualys:
  username: "qualys_username"
  password: "qualys_password"
  option: "Name Of The Qualys Scan Profile"
  proxy_url: "http://in.case.you.need.to.proxy.to.reach.qualys/"
  url: "https://qualysapi.qualys.com/api/2.0/fo/scan/"
  drop6: true
```

Testing
-------
There is one type of test file.  The `*_test.go` are standard golang unit test
files.  The examples can be run as integration tests.

License
-------
Apache v2.

Contributing
------------
The code repository utilizes the OpenStack CI infrastructure.  Please use the
[recommended
workflow](http://docs.openstack.org/infra/manual/developers.html#development-workflow).
If you are not a member yet, please consider joining as an [OpenStack
contributor](http://docs.openstack.org/infra/manual/developers.html).  If you
have questions or comments, you can email the maintainer(s).

Coding Style
------------
The source code is automatically formatted to follow `go fmt`.

OpenStack Environment
---------------------
* Release note management is done using [reno](https://docs.openstack.org/reno/latest/user/usage.html)
* Zuul CI jobs are defined in-repo, [using these techniques](https://docs.openstack.org/infra/manual/zuulv3.html#howto-in-repo)
