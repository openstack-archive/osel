package main

import "log"

type SyslogTestActions struct {
	savedLogs []string
}

func (s *SyslogTestActions) Connect() error {
	return nil
}

func (s *SyslogTestActions) Info(writeMe string) {
	log.Printf("FAKE SYSLOG LINE: %s\n", writeMe)
	s.savedLogs = append(s.savedLogs, writeMe)
}

func (s *SyslogTestActions) GetLogs() []string {
	return s.savedLogs
}

func connectFakeSyslog() *SyslogTestActions {
	return new(SyslogTestActions)
}
