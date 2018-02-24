package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// OselVersion is exposed in the logged JSON in the "source_type" field.  This
// will allow us to track the version of the logging specification.
// 1.0: Initial revision
// 1.1: Added qualys_scan_id and qualys_scan_error
const OselVersion = "osel1.1"

// Debug is a global variable to toggle debug logging
var Debug bool

// RabbitMQ URI
var rabbitURI string

func main() {
	// Declare the configuration
	viperConfigs := []ViperConfig{
		ViperConfig{Key: "batch_interval", Description: "Interval of time in minutes for message batching"},
		ViperConfig{Key: "debug", Description: "Output additional messages for debugging"},
		ViperConfig{Key: "rabbit_uri", Description: "AMQP connection uri.  See: https://www.rabbitmq.com/uri-spec.html"},
		ViperConfig{Key: "openstack.identity_endpoint", Description: "Openstack Keystone Endpoint"},
		ViperConfig{Key: "openstack.user", Description: "Openstack user that has at least read only access to all tenants/ports/security groups in the region."},
		ViperConfig{Key: "openstack.password", Description: "Password for the Openstack user"},
		ViperConfig{Key: "openstack.region", Description: "The name of the region running this process"},
		ViperConfig{Key: "qualys.drop6", Description: "Should IPv6 addresses be incorporated in Qualys scans?  true or false."},
		ViperConfig{Key: "qualys.username", Description: "Username for credentials for the Qualys external scanning service"},
		ViperConfig{Key: "qualys.password", Description: "Password for credentials for the Qualys external scanning service"},
		ViperConfig{Key: "qualys.url", Description: "URL for thw Qualys service"},
		ViperConfig{Key: "qualys.proxy_url", Description: "URL for an HTTP proxy that will permit access to the Qualys service"},
		ViperConfig{Key: "syslog_server", Description: "FQDN of the server for events to log to over the network"},
		ViperConfig{Key: "syslog_port", Description: "Port for communication to syslog, defaults to 514"},
		ViperConfig{Key: "syslog_protocol", Description: "tcp or udp, defaults to tcp"},
		ViperConfig{Key: "retry_syslog", Description: "Should the process keep trying if it cannot reach syslog?  true or false."},
	}
	configPath := os.Getenv("EL_CONFIG") //The config path comes from ENV.
	if configPath == "" {
		log.Fatalln("Fatal Error: The Config file was not set to EL_CONFIG.")
	}
	if err := InitViper(configPath, viperConfigs); err != nil {
		log.Fatalf("Fatal Error: (%s) while reading config file %s", err, configPath)
	}

	// Set defaults
	viper.SetDefault("batch_interval", 60)
	viper.SetDefault("debug", true)
	viper.SetDefault("qualys.drop6", true)
	viper.SetDefault("qualys.url", "https://qualysapi.qualys.com/api/2.0/fo/scan/")
	viper.SetDefault("syslog_port", "514")
	viper.SetDefault("syslog_protocol", "tcp")

	// Watch for config changes
	viper.WatchConfig()
	viper.OnConfigChange(func(fsnotify.Event) {
		if err := ValidateConfig(viperConfigs); err != nil {
			log.Printf("Fatal Error: %s while refreshing config file %s\n", err, configPath)
		}
	})

	batchInterval := viper.GetInt("batch_interval")
	Debug = viper.GetBool("debug")

	// Initialize AMQP
	rabbitURI = viper.GetString("rabbit_uri")
	amqpBus := new(AmqpActions)
	amqpBus.Options = AmqpOptions{
		RabbitURI: rabbitURI,
	}
	amqpIncoming, amqpErrorNotify, err := amqpBus.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize Qualys
	qualysURL, err := url.Parse(viper.GetString("qualys.url"))
	if err != nil {
		log.Fatal(err)
	}
	qualysProxyURL, err := url.Parse(viper.GetString("qualys.proxy_url"))
	if err != nil {
		log.Fatal(err)
	}
	qualys := new(QualysActions)
	qualys.Options = QualysOptions{
		DropIPv6:       viper.GetBool("qualys.drop6"),
		Password:       viper.GetString("qualys.password"),
		ProxyURL:       qualysProxyURL,
		QualysURL:      qualysURL,
		ScanOptionName: viper.GetString("qualys.option"),
		MinRemaining:   viper.GetInt("qualys.min_remaining"),
		UserName:       viper.GetString("qualys.username"),
	}

	// Initialize OpenStack
	openstack := new(OpenStackActions)
	openstack.Options = OpenStackOptions{
		KeystoneURI: viper.GetString("openstack.identity_endpoint"),
		Password:    viper.GetString("openstack.password"),
		RegionName:  viper.GetString("openstack.region"),
		UserName:    viper.GetString("openstack.user"),
	}

	// Initialize Syslog
	logger := new(SyslogActions)
	logger.Options = SyslogOptions{
		Host:     viper.GetString("syslog_server"),
		Port:     viper.GetString("syslog_port"),
		Protocol: viper.GetString("syslog_protocol"),
		Retry:    viper.GetBool("retry_syslog"),
	}
	err = logger.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// run main loop
	batchDuration := time.Duration(batchInterval) * time.Minute
	mainLoop(batchDuration, amqpIncoming, amqpErrorNotify, openstack, logger, qualys)
	defer amqpBus.Close()
}
