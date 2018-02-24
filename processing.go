package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/streadway/amqp"
)

func processWaitingEvent(delivery amqp.Delivery, openstackActions OpenStackActioner) (Event, error) {
	// executes when an event is waiting
	event, err := ParseEvent(delivery.Body)
	if err != nil {
		return Event{}, fmt.Errorf("Failed to parse event due to error: %s", err)
	}
	if event.Processor == nil {
		if !Debug {
			return Event{}, nil
		}
		return Event{}, fmt.Errorf("Ignoring event type %s", event.EventData.EventType)
	}
	if Debug {
		log.Printf("Processing event type %s\n", event.EventData.EventType)
	}
	err = event.Processor.FillExtraData(&event, openstackActions)
	if err != nil {
		return Event{}, fmt.Errorf("Error fetching extra data: %s", err)
	}
	return event, nil
}

func logEvents(events []Event, logger SyslogActioner, qualys QualysActioner) {
	var ipAddresses []string
	var qualysIPAddresses []string

	if Debug {
		log.Println("Timer Expired")
	}

	// De-dupe IP addresses and get them into a single struct
	dedupIPAddresses := make(map[string]struct{})
	for _, event := range events {
		for _, IPs := range event.IPs {
			for _, IP := range IPs {
				if _, ok := dedupIPAddresses[IP]; !ok {
					ipAddresses = append(ipAddresses, IP)
				}
				dedupIPAddresses[IP] = struct{}{}
			}
		}
	}

	// Disregard the scan if no targets have been found
	if len(ipAddresses) == 0 {
		if Debug {
			log.Println("Nothing to scan, skipping...")
		}
		return
	}

	// Remove IPv6 addresses
	if qualys.DropIPv6() {
		for ipAddressIndex := range ipAddresses {
			testIPAddress := ipAddresses[ipAddressIndex]
			if net.ParseIP(testIPAddress).To4() != nil {
				qualysIPAddresses = append(qualysIPAddresses, testIPAddress)
			} else {
				log.Println("Disregarded IPv6 address", testIPAddress)
			}
		}
	}

	// Execute Qualys scan

	log.Println("Qualys Scan Starting")
	scanID, scanError := qualys.InitiateScan(qualysIPAddresses)
	log.Printf("Qualys Scan Complete: scan ID='%s'; scan_error='%v'", scanID, scanError)

	// Iterate through entries and format the logs
	log.Printf("Processing %d events\n", len(events))
	for _, event := range events {
		event.QualysScanID = scanID
		if scanError != nil {
			event.QualysScanError = scanError.Error()
		}
		event.LogLines, _ = event.Processor.FormatLogs(&event, qualysIPAddresses)

		// Output the logs
		log.Printf("Processing %d loglines\n", len(event.LogLines))
		for lineToLog := range event.LogLines {
			logger.Info(event.LogLines[lineToLog])
		}
	}
}

func mainLoop(batchInterval time.Duration, deliveries <-chan amqp.Delivery, amqpNotifyError chan *amqp.Error, openstackActions OpenStackActioner, logger SyslogActioner, qualys QualysActioner) {
	var events []Event
	ticker := time.NewTicker(batchInterval)
	amqpReconnectTimer := time.NewTimer(1)
	for {
		select {
		case e := <-deliveries:
			event, err := processWaitingEvent(e, openstackActions)
			if err != nil {
				log.Printf("Event skipped: %s\n", err)
				continue
			}
			events = append(events, event)
		case <-ticker.C:
			logEvents(events, logger, qualys)
			events = nil
		case err := <-amqpNotifyError:
			// Reinitialize AMQP on connection error
			log.Printf("AMQP connection error: %s\n", err)
			amqpReconnectTimer = time.NewTimer(time.Second * 30)
		case <-amqpReconnectTimer.C:
			var err error
			amqpBus := new(AmqpActions)
			amqpBus.Options = AmqpOptions{
				RabbitURI: rabbitURI,
			}
			deliveries, amqpNotifyError, err = amqpBus.Connect()
			if err != nil {
				log.Printf("AMQP retry connection error: %s\n", err)
				amqpReconnectTimer = time.NewTimer(time.Second * 30)
			} else {
				log.Printf("AMQP reconnected\n")
			}
		}
	}
}
