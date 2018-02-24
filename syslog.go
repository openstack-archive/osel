package main

/*

syslog - This file includes all of the logic necessary to interact with syslog.
This is extrapolated out so that a SyslogActioner interface can be
passed to functions.  Doing this allows testing by mock classes to be created
that can be passed to functions.

Since this is a wrapper around the log/syslog library, this does not need
testing.

*/

import (
	"fmt"
	"log"
	"log/syslog"
	"net"
)

// SyslogActioner is an interface for an SyslogActions class.  Having
// this as an interface allows us to pass in a dummy class for testing that
// just returns mocked data.
type SyslogActioner interface {
	Connect() error
	Info(string)
}

// SyslogActions is a class that handles all interactions directly with Syslog.
// See the comment on SyslogActioner for rationale.
type SyslogActions struct {
	logger  *syslog.Writer
	Options SyslogOptions
}

// SyslogOptions is a class to convey all of the configurable options for the
// SyslogActions class.
type SyslogOptions struct {
	Host     string
	Protocol string
	Retry    bool
	Port     string
}

// Info is the main method for the SyslogActioner class, it writes an
// info-level message to the syslog stream.
func (s SyslogActions) Info(writeMe string) {
	log.Println("Logged: ", writeMe)
	s.logger.Info(writeMe)
}

// Connect is the method that establishes the connection to the syslog server
// over the network.
func (s *SyslogActions) Connect() error {
	var err error

	address := net.JoinHostPort(s.Options.Host, s.Options.Port)

	if Debug {
		log.Printf("Opening %q syslog socket to %q\n", s.Options.Protocol, s.Options.Host)
	}

	s.logger, err = syslog.Dial(s.Options.Protocol, address, syslog.LOG_INFO, "osel")
	if err != nil {
		log.Printf("error opening syslog socket to %s: %s\n", s.Options.Host, err)
		if s.Options.Retry {
			for err != nil {
				log.Println("retrying")
				s.logger, err = syslog.Dial(s.Options.Protocol, address, syslog.LOG_INFO, "osel")
			}
		}
		return fmt.Errorf("error opening syslog socket to %s: %s", s.Options.Host, err)
	}

	if Debug {
		log.Println("Successfully connected to syslog host", s.Options.Host)
	}

	return nil
}
