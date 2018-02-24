package main

import (
	"log"
	//"github.com/satori/go.uuid"
)

// QualysActions is a class that handles all interactions directly with Qualys.
// See the comment on QualysActioner for rationale.
type QualysTestActions struct {
	testUUID string
}

// InitiateScan is the main method for the QualysActioner class, it
// makes a call to the Qualys API to start a scan and harvests a scan ID, and
// an optional error string if there is a problem contacting Qualys.
func (s *QualysTestActions) InitiateScan(ipAddresses []string) (string, error) {
	//testUUID = uuid.NewV4().String()
	s.testUUID = `5fbf3cef-976e-475d-bd84-47ef23638a6b`
	log.Printf("FAKE QUALYS SCAN: %s\n", s.testUUID)
	return s.testUUID, nil
}

// GetTestScanID returns the fake UUID created in testing. This allows for
// inspection of the UUID in unit tests.
func (s *QualysTestActions) GetTestScanID() string {
	return s.testUUID
}

func (s *QualysTestActions) DropIPv6() bool {
	return false
}

func connectFakeQualys() *QualysTestActions {
	return new(QualysTestActions)
}
