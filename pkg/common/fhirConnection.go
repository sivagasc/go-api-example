package common

import (
	"sync"

	"github.com/sivagasc/go-api-example/pkg/fhir"
)

//FHIRConnection variable stores the FHIR Connection
var FHIRConnection *fhir.Connection

var fhirOnce sync.Once

//ConnectToFHIR Method connect to FHIR server
func ConnectToFHIR(serverURL string) *fhir.Connection {

	fhirOnce.Do(func() {
		logger.Info().Msg("Connecting to FHIR....")
		FHIRConnection = fhir.NewConnection(serverURL)
	})
	logger.Info().Msg("FHIR Connection success:")

	return FHIRConnection

}

//GetFHIRConnection Get the Database connection
func GetFHIRConnection() *fhir.Connection {
	logger.Debug().Msg("Get FHIR Connection..")
	return FHIRConnection
}
