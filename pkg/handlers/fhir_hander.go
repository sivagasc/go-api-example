package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sivagasc/go-api-example/pkg/common"
	"github.com/sivagasc/go-api-example/pkg/fhir"
	"github.com/sivagasc/go-api-example/pkg/utils"
)

// GetPatient is a handler method to return Patient details
func GetPatient() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Get Logger
		logger := common.GetLoggerInstance()

		logger.Info().Msg("*** Get Patient details")

		id := mux.Vars(req)["id"]
		if id == "" {
			utils.RespondError(w, http.StatusBadRequest, "Expected id as an input.")
			logger.Error().Msg("Expected id as an input.")
			return
		}

		fhirConnection := fhir.NewConnection("http://test.fhir.org/r4/")
		logger.Info().Msg("FHIR Connection success:")

		// Get Patient details
		pat, err := fhirConnection.GetPatient(id)
		if err != nil {
			fmt.Println("Error", err.Error())
			utils.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondJSON(w, http.StatusOK, pat)
		return
	})
}
